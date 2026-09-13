// Package queue maps Golib queue jobs to and from CloudEvents while retaining
// execution, retry, settlement, and operational state outside the event.
package queue

import (
	"fmt"
	"time"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	"github.com/faustbrian/go-queue/job"
	"github.com/faustbrian/go-tenancy"
)

var (
	// ErrInvalidInput reports a queue message, option, tenant, extension, or
	// CloudEvent that cannot satisfy the adapter contract. Its zero value is
	// not meaningful; callers should compare returned errors with errors.Is.
	ErrInvalidInput = cloudevents.ErrInvalidAdapterInput
	// ErrMetadataCollision reports portable CloudEvent metadata that conflicts
	// with caller-retained queue state. The retained value is never overwritten.
	ErrMetadataCollision = cloudevents.ErrMetadataCollision
)

// Loss describes one CloudEvent field that the queue model cannot represent.
// The zero value describes no useful loss; values are owned by their Report.
type Loss = cloudevents.AdapterLoss

// Report lists deliberate information loss during a conversion. Its zero
// value means the conversion was lossless, and callers own the returned value.
type Report = cloudevents.AdapterReport

// Options supplies CloudEvent identity and occurrence metadata that a queue
// message does not necessarily own. The zero value is invalid because Source
// and, absent retained metadata, StableID and Type are required.
type Options struct {
	// Source is the required CloudEvent source URI-reference. It is caller-owned.
	Source string
	// StableID overrides Metadata.OriginalID. An empty value uses the retained
	// original ID; conversion fails when both are empty.
	StableID string
	// Type overrides Metadata.JobType. An empty value uses the retained job type;
	// conversion fails when both are empty.
	Type string
	// OccurredAt is optional CloudEvent occurrence time. Nil omits the attribute;
	// a non-nil value is copied and is not retained by reference.
	OccurredAt *time.Time
}

// ToCloudEvent converts message into a CloudEvent and returns a deep copy of
// the queue-owned state needed for reversal. Portable identity and context are
// copied into the event, collisions or invalid values fail, and neither the
// input message nor caller-owned slices, maps, or timestamps are retained.
func ToCloudEvent(message job.Message, options Options) (cloudevents.Event, job.Message, Report, error) {
	if message.Metadata != nil {
		if _, err := validatedTenant(message.Metadata.TenantID); err != nil {
			return cloudevents.Event{}, job.Message{}, Report{}, err
		}
	}
	if err := message.Validate(); err != nil || options.Source == "" {
		return cloudevents.Event{}, job.Message{}, Report{}, fmt.Errorf("%w: queue mapping", ErrInvalidInput)
	}
	id, eventType, contentType := options.StableID, options.Type, ""
	if message.Metadata != nil {
		if id == "" {
			id = message.Metadata.OriginalID
		}
		if eventType == "" {
			eventType = message.Metadata.JobType
		}
		contentType = message.Metadata.ContentType
	}
	if id == "" || eventType == "" {
		return cloudevents.Event{}, job.Message{}, Report{}, fmt.Errorf("%w: queue identity", ErrInvalidInput)
	}
	data, err := adapter.DataFromPayload(contentType, message.Body)
	if err != nil {
		return cloudevents.Event{}, job.Message{}, Report{}, err
	}
	extensions, err := queueExtensions(message.Metadata)
	if err != nil {
		return cloudevents.Event{}, job.Message{}, Report{}, err
	}
	event, err := cloudevents.NewEvent(cloudevents.Attributes{ID: id, Source: options.Source, Type: eventType, DataContentType: contentType, Time: adapter.CloneTimePointer(options.OccurredAt), Extensions: extensions}, data)
	if err != nil {
		return cloudevents.Event{}, job.Message{}, Report{}, err
	}
	return event, cloneJob(message), Report{}, nil
}

// FromCloudEvent restores portable data into a deep copy of state. The caller
// retains ownership of queue-only execution and settlement fields; conflicting
// retained portable metadata returns ErrMetadataCollision instead of choosing
// either value. The returned Report records CloudEvent fields queue cannot own.
func FromCloudEvent(event cloudevents.Event, state job.Message) (job.Message, Report, error) {
	if err := event.Validate(); err != nil || !event.Data().Present() {
		return job.Message{}, Report{}, fmt.Errorf("%w: queue target", ErrInvalidInput)
	}
	state = cloneJob(state)
	if state.Metadata == nil {
		state.Metadata = &job.Metadata{}
	}
	if state.Metadata.OriginalID != "" && state.Metadata.OriginalID != event.ID() {
		return job.Message{}, Report{}, fmt.Errorf("%w: queue original id", ErrMetadataCollision)
	}
	if state.Metadata.JobType != "" && state.Metadata.JobType != event.Type() {
		return job.Message{}, Report{}, fmt.Errorf("%w: queue job type", ErrMetadataCollision)
	}
	state.Body = adapter.RestoreRetainedNil(event.Data(), state.Body == nil)
	state.Metadata.OriginalID = event.ID()
	state.Metadata.JobType = event.Type()
	contentType, _ := event.DataContentType()
	if state.Metadata.ContentType != "" && state.Metadata.ContentType != contentType {
		return job.Message{}, Report{}, fmt.Errorf("%w: queue content type", ErrMetadataCollision)
	}
	state.Metadata.ContentType = contentType
	if err := verifyExtensions(event, state.Metadata); err != nil {
		return job.Message{}, Report{}, err
	}
	report := Report{Losses: []Loss{{Field: "source", Reason: "not represented by queue"}}}
	adapter.AppendDataKindLoss(event, &report, "queue", true)
	adapter.AppendOptionalContextLosses(event, &report, "queue")
	adapter.AppendExtensionLosses(event, &report, "queue", map[string]struct{}{adapter.CorrelationIDExtension: {}, adapter.RequestIDExtension: {}, adapter.CausationIDExtension: {}, "traceparent": {}, "tracestate": {}, adapter.TenantIDExtension: {}})
	return state, report, nil
}

func cloneJob(message job.Message) job.Message {
	message.Body = adapter.CloneBytes(message.Body)
	if message.Metadata == nil {
		return message
	}
	metadata := *message.Metadata
	metadata.Tags = adapter.CloneStrings(metadata.Tags)
	metadata.Correlation = adapter.CloneStrings(metadata.Correlation)
	metadata.TraceContext = adapter.CloneStrings(metadata.TraceContext)
	if metadata.EnqueuedAt != nil {
		enqueuedAt := *metadata.EnqueuedAt
		metadata.EnqueuedAt = &enqueuedAt
	}
	message.Metadata = &metadata
	return message
}

func queueExtensions(metadata *job.Metadata) (map[string]cloudevents.Attribute, error) {
	extensions := map[string]cloudevents.Attribute{}
	if metadata == nil {
		return extensions, nil
	}
	if _, err := validatedTenant(metadata.TenantID); err != nil {
		return nil, err
	}
	for _, mapping := range []struct{ target, source string }{
		{adapter.CorrelationIDExtension, metadata.Correlation[adapter.CorrelationIDExtension]},
		{adapter.RequestIDExtension, metadata.Correlation[adapter.RequestIDExtension]},
		{adapter.CausationIDExtension, metadata.Correlation[adapter.CausationIDExtension]},
		{"traceparent", metadata.TraceContext["traceparent"]}, {"tracestate", metadata.TraceContext["tracestate"]},
		{adapter.TenantIDExtension, metadata.TenantID},
	} {
		if mapping.source != "" {
			attribute, err := cloudevents.NewStringAttribute(mapping.source)
			if err != nil {
				return nil, fmt.Errorf("%w: queue extension %s: %w", ErrInvalidInput, mapping.target, err)
			}
			extensions[mapping.target] = attribute
		}
	}
	return extensions, nil
}

func verifyExtensions(event cloudevents.Event, metadata *job.Metadata) error {
	if _, err := validatedTenant(metadata.TenantID); err != nil {
		return err
	}
	expected := []struct{ name, retained string }{
		{adapter.CorrelationIDExtension, metadata.Correlation[adapter.CorrelationIDExtension]},
		{adapter.RequestIDExtension, metadata.Correlation[adapter.RequestIDExtension]},
		{adapter.CausationIDExtension, metadata.Correlation[adapter.CausationIDExtension]},
		{"traceparent", metadata.TraceContext["traceparent"]}, {"tracestate", metadata.TraceContext["tracestate"]},
		{adapter.TenantIDExtension, metadata.TenantID},
	}
	for _, mapping := range expected {
		value, present, err := adapter.StringExtension(event, mapping.name)
		if err != nil {
			return err
		}
		if mapping.name == adapter.TenantIDExtension && present {
			if _, err := validatedTenant(value); err != nil {
				return err
			}
		}
		if (!present && mapping.retained != "") || (present && mapping.retained != value) {
			return fmt.Errorf("%w: queue extension %s", ErrMetadataCollision, mapping.name)
		}
	}
	return nil
}

func validatedTenant(value string) (string, error) {
	if value == "" {
		return "", nil
	}
	tenantID, err := tenancy.ParseTenantID(value)
	if err != nil {
		return "", fmt.Errorf("%w: tenant: %w", ErrInvalidInput, err)
	}
	return tenantID.Value(), nil
}
