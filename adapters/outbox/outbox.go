// Package outbox maps Golib transactional-outbox envelopes to and from
// CloudEvents while retaining relay-owned state outside the event.
package outbox

import (
	"fmt"
	"time"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	goliboutbox "github.com/faustbrian/go-transactional-outbox"
)

var (
	// ErrInvalidInput reports an invalid envelope, retained state, payload, or
	// CloudEvent. Callers can match it with errors.Is.
	ErrInvalidInput = cloudevents.ErrInvalidAdapterInput
	// ErrMetadataCollision reports that a retained non-empty outbox identifier
	// disagrees with the CloudEvent identifier. Callers can match it with errors.Is.
	ErrMetadataCollision = cloudevents.ErrMetadataCollision
)

// Loss describes one CloudEvent field that an outbox envelope cannot
// reconstruct. It aliases the root module's stable loss contract.
type Loss = cloudevents.AdapterLoss

// Report lists intentional, field-specific conversion losses. Its zero value
// means no loss was observed.
type Report = cloudevents.AdapterReport

// Options supplies CloudEvent context that is not owned by an outbox envelope.
// Its zero value is invalid because Source and Type are required.
type Options struct {
	// Source identifies the producer and must be a valid CloudEvents URI-reference.
	Source string
	// Type names the portable event; it is intentionally distinct from the retained Topic.
	Type string
	// DataContentType optionally describes Payload and controls JSON validation.
	DataContentType string
	// DataSchema optionally identifies the payload schema.
	DataSchema string
	// Subject optionally supplies portable subject context not owned by the outbox.
	Subject string
	// OccurredAt optionally supplies occurrence time. The pointed-to value is cloned.
	OccurredAt *time.Time
}

// ToCloudEvent maps envelope into a CloudEvent and returns a deep-cloned copy
// as caller-owned retained state. It rejects invalid identifiers, attributes,
// and declared payload encodings; successful conversion has no declared losses.
func ToCloudEvent(envelope goliboutbox.Envelope, options Options) (cloudevents.Event, goliboutbox.Envelope, Report, error) {
	if envelope.ID == "" || options.Source == "" || options.Type == "" {
		return cloudevents.Event{}, goliboutbox.Envelope{}, Report{}, fmt.Errorf("%w: outbox mapping", ErrInvalidInput)
	}
	data, err := adapter.DataFromPayload(options.DataContentType, envelope.Payload)
	if err != nil {
		return cloudevents.Event{}, goliboutbox.Envelope{}, Report{}, err
	}
	event, err := cloudevents.NewEvent(cloudevents.Attributes{ID: envelope.ID, Source: options.Source, Type: options.Type, DataContentType: options.DataContentType, DataSchema: options.DataSchema, Subject: options.Subject, Time: adapter.CloneTimePointer(options.OccurredAt)}, data)
	if err != nil {
		return cloudevents.Event{}, goliboutbox.Envelope{}, Report{}, err
	}
	return event, cloneEnvelope(envelope), Report{}, nil
}

// FromCloudEvent reconstructs an outbox envelope from event and caller-owned
// state. It preserves nil versus non-nil empty payloads, clones retained maps
// and bytes, rejects identifier collisions, and reports every portable field
// the outbox model cannot represent.
func FromCloudEvent(event cloudevents.Event, state goliboutbox.Envelope) (goliboutbox.Envelope, Report, error) {
	if err := event.Validate(); err != nil || state.Topic == "" || state.PayloadVersion == 0 {
		return goliboutbox.Envelope{}, Report{}, fmt.Errorf("%w: outbox target", ErrInvalidInput)
	}
	if state.ID != "" && state.ID != event.ID() {
		return goliboutbox.Envelope{}, Report{}, fmt.Errorf("%w: outbox id", ErrMetadataCollision)
	}
	if !event.Data().Present() {
		return goliboutbox.Envelope{}, Report{}, fmt.Errorf("%w: outbox payload", ErrInvalidInput)
	}
	state = cloneEnvelope(state)
	state.ID = event.ID()
	state.Payload = adapter.RestoreRetainedNil(event.Data(), state.Payload == nil)
	report := Report{Losses: []Loss{{Field: "source", Reason: "not represented by outbox"}, {Field: "type", Reason: "outbox topic is transport-owned"}}}
	if _, present := event.DataContentType(); present {
		report.Losses = append(report.Losses, Loss{Field: "datacontenttype", Reason: "not represented by outbox"})
	}
	adapter.AppendDataKindLoss(event, &report, "outbox", false)
	adapter.AppendOptionalContextLosses(event, &report, "outbox")
	adapter.AppendExtensionLosses(event, &report, "outbox", nil)
	return state, report, nil
}

func cloneEnvelope(envelope goliboutbox.Envelope) goliboutbox.Envelope {
	envelope.Payload = adapter.CloneBytes(envelope.Payload)
	envelope.Metadata = adapter.CloneStrings(envelope.Metadata)
	return envelope
}
