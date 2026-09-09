// Package eventsourcing maps persisted Golib event-sourcing messages to and
// from CloudEvents while keeping event-store state explicit and caller-owned.
package eventsourcing

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"time"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	golibeventsourcing "github.com/faustbrian/go-event-sourcing"
	"github.com/faustbrian/go-tenancy"
)

const eventSchemaExtension = "eventschema"

var (
	// ErrInvalidInput reports an invalid source message, retained state, tenant,
	// payload, or CloudEvent. Callers can match it with errors.Is.
	ErrInvalidInput = cloudevents.ErrInvalidAdapterInput
	// ErrMetadataCollision reports that trusted retained state disagrees with
	// portable CloudEvent metadata. Callers can match it with errors.Is.
	ErrMetadataCollision = cloudevents.ErrMetadataCollision
	// ErrUntrustedMetadata reports that reconstruction would require accepting
	// CloudEvent metadata which the caller has not marked as trusted.
	ErrUntrustedMetadata = cloudevents.ErrUntrustedMetadata
)

// Loss describes one CloudEvent field that the event-sourcing message cannot
// reconstruct. It aliases the root module's stable loss contract.
type Loss = cloudevents.AdapterLoss

// Report lists intentional, field-specific conversion losses. Its zero value
// means no loss was observed.
type Report = cloudevents.AdapterReport

// Options supplies CloudEvent context that is not owned by an event-store
// message. Its zero value is invalid because Source is required.
type Options struct {
	// Source identifies the event producer and must be a valid CloudEvents URI-reference.
	Source string
	// DataSchema optionally identifies the schema of the encoded event payload.
	DataSchema string
	// OccurredAt optionally supplies domain occurrence time; it is cloned and is
	// intentionally distinct from the retained RecordedAt persistence time.
	OccurredAt *time.Time
}

// State retains event-store-owned fields that CloudEvents cannot represent
// without changing their meaning. Its zero value is invalid for reconstruction.
// Callers own this value and must store it alongside the CloudEvent.
type State struct {
	// Stream is the original aggregate stream and is bound to the CloudEvent subject.
	Stream golibeventsourcing.StreamID
	// StreamVersion is the positive optimistic-concurrency version to reconstruct.
	StreamVersion uint64
	// GlobalPosition is the optional store-wide ordering position; zero means absent.
	GlobalPosition golibeventsourcing.GlobalPosition
	// EventVersion is the positive encoded-event schema version and is collision-checked.
	EventVersion golibeventsourcing.SchemaVersion
	// Metadata is caller-owned event-store metadata. Conversion clones the map.
	Metadata map[string]string
	// RecordedAt is the required persistence timestamp, not CloudEvents occurrence time.
	RecordedAt time.Time
	// CorrelationID is retained and collision-checked against the portable extension.
	CorrelationID string
	// CausationID is retained and collision-checked against the portable extension.
	CausationID string
	// Tenant is retained, validated, and collision-checked against the portable extension.
	Tenant string
	// Partition is retained and collision-checked against the portable extension.
	Partition string
	// TrustMetadata permits missing retained correlation, causation, tenant, or
	// partition values to be sourced from their CloudEvent extensions.
	TrustMetadata bool
}

// ToCloudEvent maps message into a CloudEvent and returns independently owned
// retained state. It rejects invalid payloads, tenants, identifiers, and event
// attributes; successful conversion has no declared losses.
func ToCloudEvent(message golibeventsourcing.Message, options Options) (cloudevents.Event, State, Report, error) {
	if options.Source == "" || message.ID().IsZero() {
		return cloudevents.Event{}, State{}, Report{}, fmt.Errorf("%w: event-sourcing mapping", ErrInvalidInput)
	}
	encoded := message.Event()
	data, err := adapter.DataFromPayload(encoded.ContentType(), encoded.Payload())
	if err != nil {
		return cloudevents.Event{}, State{}, Report{}, err
	}
	extensions := map[string]cloudevents.Attribute{}
	extensions[eventSchemaExtension], _ = cloudevents.NewStringAttribute(strconv.FormatUint(uint64(encoded.Version()), 10))
	state := State{Stream: message.Stream(), StreamVersion: message.StreamVersion(), EventVersion: encoded.Version(), Metadata: message.Metadata(), RecordedAt: message.RecordedAt(), TrustMetadata: true}
	if position, present := message.GlobalPosition(); present {
		state.GlobalPosition = position
	}
	if value, present := message.CorrelationID(); present {
		state.CorrelationID = value.String()
		extensions[adapter.CorrelationIDExtension], _ = cloudevents.NewStringAttribute(state.CorrelationID)
	}
	if value, present := message.CausationID(); present {
		state.CausationID = value.String()
		extensions[adapter.CausationIDExtension], _ = cloudevents.NewStringAttribute(state.CausationID)
	}
	if value, present := message.Tenant(); present {
		state.Tenant, err = validatedTenant(value)
		if err != nil {
			return cloudevents.Event{}, State{}, Report{}, err
		}
		extensions[adapter.TenantIDExtension], _ = cloudevents.NewStringAttribute(state.Tenant)
	}
	if value, present := message.Partition(); present {
		state.Partition = value
		extensions["partitionkey"], _ = cloudevents.NewStringAttribute(value)
	}
	state.Metadata = adapter.CloneStrings(state.Metadata)
	event, err := cloudevents.NewEvent(cloudevents.Attributes{ID: message.ID().String(), Source: options.Source, Type: encoded.Name().String(), DataContentType: encoded.ContentType(), DataSchema: options.DataSchema, Subject: encodeStreamSubject(message.Stream()), Time: adapter.CloneTimePointer(options.OccurredAt), Extensions: extensions}, data)
	if err != nil {
		return cloudevents.Event{}, State{}, Report{}, err
	}
	return event, state, Report{}, nil
}

// FromCloudEvent reconstructs an event-store message from event and caller-owned
// state. Retained metadata is cloned; trusted fields are collision-checked, and
// every portable field that cannot be reconstructed is reported as a loss.
func FromCloudEvent(event cloudevents.Event, state State) (golibeventsourcing.Message, Report, error) {
	if err := event.Validate(); err != nil || state.Stream.IsZero() || state.StreamVersion == 0 || state.EventVersion == 0 || state.RecordedAt.IsZero() {
		return golibeventsourcing.Message{}, Report{}, fmt.Errorf("%w: event-sourcing target", ErrInvalidInput)
	}
	if subject, present := event.Subject(); !present || subject != encodeStreamSubject(state.Stream) {
		return golibeventsourcing.Message{}, Report{}, fmt.Errorf("%w: subject", ErrMetadataCollision)
	}
	version, _, _ := adapter.StringExtension(event, eventSchemaExtension)
	if version != strconv.FormatUint(uint64(state.EventVersion), 10) {
		return golibeventsourcing.Message{}, Report{}, fmt.Errorf("%w: eventschema", ErrMetadataCollision)
	}
	correlationID, err := adapter.MappedString(event, adapter.CorrelationIDExtension, state.CorrelationID, state.TrustMetadata)
	if err != nil {
		return golibeventsourcing.Message{}, Report{}, err
	}
	causationID, err := adapter.MappedString(event, adapter.CausationIDExtension, state.CausationID, state.TrustMetadata)
	if err != nil {
		return golibeventsourcing.Message{}, Report{}, err
	}
	tenant, err := adapter.MappedString(event, adapter.TenantIDExtension, state.Tenant, state.TrustMetadata)
	if err != nil {
		return golibeventsourcing.Message{}, Report{}, err
	}
	tenant, err = validatedTenant(tenant)
	if err != nil {
		return golibeventsourcing.Message{}, Report{}, err
	}
	partition, err := adapter.MappedString(event, "partitionkey", state.Partition, state.TrustMetadata)
	if err != nil {
		return golibeventsourcing.Message{}, Report{}, err
	}
	contentType, present := event.DataContentType()
	if !present || !event.Data().Present() {
		return golibeventsourcing.Message{}, Report{}, fmt.Errorf("%w: event data", ErrInvalidInput)
	}
	encoded, err := golibeventsourcing.NewEncodedEvent(golibeventsourcing.EncodedEventInput{Name: event.Type(), Version: state.EventVersion, ContentType: contentType, Payload: event.Data().Bytes()})
	if err != nil {
		return golibeventsourcing.Message{}, Report{}, err
	}
	pending, err := golibeventsourcing.NewPendingMessage(golibeventsourcing.PendingMessageInput{ID: event.ID(), Stream: state.Stream, Event: encoded, Metadata: adapter.CloneStrings(state.Metadata), RecordedAt: state.RecordedAt, CorrelationID: correlationID, CausationID: causationID, Tenant: tenant, Partition: partition})
	if err != nil {
		return golibeventsourcing.Message{}, Report{}, err
	}
	message, _ := golibeventsourcing.NewMessage(golibeventsourcing.MessageInput{Pending: pending, StreamVersion: state.StreamVersion, GlobalPosition: state.GlobalPosition})
	report := Report{Losses: []Loss{{Field: "source", Reason: "not represented by event-sourcing"}}}
	adapter.AppendDataKindLoss(event, &report, "event-sourcing", true)
	if _, present := event.DataSchema(); present {
		report.Losses = append(report.Losses, Loss{Field: "dataschema", Reason: "not represented by event-sourcing"})
	}
	if _, present := event.Time(); present {
		report.Losses = append(report.Losses, Loss{Field: "time", Reason: "recorded_at is not occurrence time"})
	}
	adapter.AppendExtensionLosses(event, &report, "event-sourcing", map[string]struct{}{eventSchemaExtension: {}, adapter.CorrelationIDExtension: {}, adapter.CausationIDExtension: {}, adapter.TenantIDExtension: {}, "partitionkey": {}})
	return message, report, nil
}

func encodeStreamSubject(stream golibeventsourcing.StreamID) string {
	return base64.RawURLEncoding.EncodeToString([]byte(stream.AggregateType())) + "." + base64.RawURLEncoding.EncodeToString([]byte(stream.AggregateID()))
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
