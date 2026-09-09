package eventsourcing

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/faustbrian/go-cloudevents"
	golibeventsourcing "github.com/faustbrian/go-event-sourcing"
	"github.com/faustbrian/go-tenancy"
)

func TestRoundTripRetainsEventStoreStateWithoutAliasing(t *testing.T) {
	t.Parallel()

	recordedAt := time.Date(2026, 8, 9, 1, 2, 3, 0, time.UTC)
	occurredAt := recordedAt.Add(-time.Minute)
	message := newMessage(t, messageSpec{
		payload: []byte(`{"order":"A/123"}`), contentType: "application/json", recordedAt: recordedAt,
		metadata: map[string]string{"owner": "domain"}, correlation: "correlation-1", causation: "cause-1",
		tenant: "tenant-a", partition: "tenant-a", globalPosition: 9,
	})
	wantPayload := append([]byte(nil), message.Event().Payload()...)
	event, state, report, err := ToCloudEvent(message, Options{
		Source: "/event-store/orders", DataSchema: "https://schemas.example/order-created-v2.json", OccurredAt: &occurredAt,
	})
	if err != nil || len(report.Losses) != 0 {
		t.Fatalf("ToCloudEvent() = %#v, %v", report, err)
	}
	message.Event().Payload()[0] = 'x'
	message.Metadata()["owner"] = "changed"
	if !bytes.Equal(event.Data().Bytes(), wantPayload) || state.Metadata["owner"] != "domain" {
		t.Fatal("conversion aliases caller-owned payload or metadata")
	}

	roundTrip, reverseReport, err := FromCloudEvent(event, state)
	if err != nil || !roundTrip.Equal(newMessage(t, messageSpec{
		payload: wantPayload, contentType: "application/json", recordedAt: recordedAt,
		metadata: map[string]string{"owner": "domain"}, correlation: "correlation-1", causation: "cause-1",
		tenant: "tenant-a", partition: "tenant-a", globalPosition: 9,
	})) {
		t.Fatalf("FromCloudEvent() = %#v, %#v, %v", roundTrip, reverseReport, err)
	}
	wantLosses := map[string]bool{"source": true, "dataschema": true, "time": true}
	for _, loss := range reverseReport.Losses {
		delete(wantLosses, loss.Field)
	}
	if len(wantLosses) != 0 {
		t.Fatalf("loss report omitted fields: %#v", wantLosses)
	}
	state.Metadata["owner"] = "retained-state-mutated"
	if roundTrip.Metadata()["owner"] != "domain" {
		t.Fatal("reconstructed message aliases caller-owned retained metadata")
	}
}

func TestToCloudEventRejectsInvalidInputs(t *testing.T) {
	t.Parallel()

	if _, _, _, err := ToCloudEvent(golibeventsourcing.Message{}, Options{Source: "/source"}); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("zero message error = %v", err)
	}
	invalidJSON := newMessage(t, messageSpec{payload: []byte("{"), contentType: "application/json"})
	if _, _, _, err := ToCloudEvent(invalidJSON, Options{Source: "/source"}); err == nil {
		t.Fatal("invalid JSON error = nil")
	}
	invalidTenant := newMessage(t, messageSpec{payload: []byte("body"), tenant: "bad?tenant"})
	if _, _, _, err := ToCloudEvent(invalidTenant, Options{Source: "/source"}); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) || !errors.Is(err, tenancy.ErrInvalidTenantID) {
		t.Fatalf("invalid tenant error = %v", err)
	}
	message := newMessage(t, messageSpec{payload: []byte("body")})
	if _, _, _, err := ToCloudEvent(message, Options{Source: "bad\nsource"}); err == nil {
		t.Fatal("invalid source error = nil")
	}
}

func TestFromCloudEventRejectsInvalidStateAndMetadataCollisions(t *testing.T) {
	t.Parallel()

	message := newMessage(t, messageSpec{
		payload: []byte("body"), metadata: map[string]string{"key": "value"}, correlation: "correlation-1",
		causation: "cause-1", tenant: "tenant-a", partition: "partition-a",
	})
	event, state, _, err := ToCloudEvent(message, Options{Source: "/source"})
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := FromCloudEvent(cloudevents.Event{}, state); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("invalid event error = %v", err)
	}

	wrongSubject := cloneEvent(t, event, func(attributes *cloudevents.Attributes) { attributes.Subject = "other" }, event.Data())
	if _, _, err := FromCloudEvent(wrongSubject, state); !errors.Is(err, cloudevents.ErrMetadataCollision) {
		t.Fatalf("subject collision error = %v", err)
	}
	withoutSchema := cloneEvent(t, event, func(attributes *cloudevents.Attributes) { delete(attributes.Extensions, eventSchemaExtension) }, event.Data())
	if _, _, err := FromCloudEvent(withoutSchema, state); !errors.Is(err, cloudevents.ErrMetadataCollision) {
		t.Fatalf("schema collision error = %v", err)
	}

	for name, mutate := range map[string]func(*State){
		"correlation": func(value *State) { value.CorrelationID = "different" },
		"causation":   func(value *State) { value.CausationID = "different" },
		"tenant":      func(value *State) { value.Tenant = "different" },
		"partition":   func(value *State) { value.Partition = "different" },
	} {
		t.Run(name, func(t *testing.T) {
			changed := state
			mutate(&changed)
			if _, _, err := FromCloudEvent(event, changed); !errors.Is(err, cloudevents.ErrMetadataCollision) {
				t.Fatalf("collision error = %v", err)
			}
		})
	}

	for name := range map[string]bool{"correlationid": true, "causationid": true, "tenantid": true, "partitionkey": true} {
		t.Run("missing "+name, func(t *testing.T) {
			missing := cloneEvent(t, event, func(attributes *cloudevents.Attributes) { delete(attributes.Extensions, name) }, event.Data())
			if _, _, err := FromCloudEvent(missing, state); !errors.Is(err, cloudevents.ErrMetadataCollision) {
				t.Fatalf("missing extension error = %v", err)
			}
		})
	}
}

func TestFromCloudEventValidatesPayloadTenantAndMessageIdentity(t *testing.T) {
	t.Parallel()

	message := newMessage(t, messageSpec{payload: []byte("body"), tenant: "tenant-a"})
	event, state, _, err := ToCloudEvent(message, Options{Source: "/source"})
	if err != nil {
		t.Fatal(err)
	}
	withoutData := cloneEvent(t, event, nil, cloudevents.Data{})
	if _, _, err := FromCloudEvent(withoutData, state); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("absent data error = %v", err)
	}
	invalidTenant, _ := cloudevents.NewStringAttribute("bad?tenant")
	invalidTenantEvent := cloneEvent(t, event, func(attributes *cloudevents.Attributes) {
		attributes.Extensions["tenantid"] = invalidTenant
	}, event.Data())
	invalidTenantState := state
	invalidTenantState.Tenant = ""
	if _, _, err := FromCloudEvent(invalidTenantEvent, invalidTenantState); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) || !errors.Is(err, tenancy.ErrInvalidTenantID) {
		t.Fatalf("invalid tenant error = %v", err)
	}
	emptyData := cloneEvent(t, event, nil, cloudevents.NewBinaryData(nil))
	if _, _, err := FromCloudEvent(emptyData, state); err == nil {
		t.Fatal("empty payload error = nil")
	}
	badID := cloneEvent(t, event, func(attributes *cloudevents.Attributes) { attributes.ID = "bad id" }, event.Data())
	if _, _, err := FromCloudEvent(badID, state); err == nil {
		t.Fatal("invalid message ID error = nil")
	}
}

func TestUntrustedPortableMetadataCanPopulateMissingRetainedFields(t *testing.T) {
	t.Parallel()

	message := newMessage(t, messageSpec{payload: []byte("body"), correlation: "correlation-1", causation: "cause-1", tenant: "tenant-a", partition: "partition-a"})
	event, state, _, err := ToCloudEvent(message, Options{Source: "/source"})
	if err != nil {
		t.Fatal(err)
	}
	state.CorrelationID, state.CausationID, state.Tenant, state.Partition = "", "", "", ""
	state.TrustMetadata = false
	if _, _, err := FromCloudEvent(event, state); !errors.Is(err, cloudevents.ErrUntrustedMetadata) {
		t.Fatalf("untrusted metadata error = %v", err)
	}
	state.TrustMetadata = true
	converted, _, err := FromCloudEvent(event, state)
	if err != nil {
		t.Fatal(err)
	}
	if value, ok := converted.Tenant(); !ok || value != "tenant-a" {
		t.Fatalf("tenant = %q, %v", value, ok)
	}
}

func TestRoundTripAllowsAbsentOptionalMetadata(t *testing.T) {
	t.Parallel()

	message := newMessage(t, messageSpec{payload: []byte("body")})
	event, state, _, err := ToCloudEvent(message, Options{Source: "/source"})
	if err != nil {
		t.Fatal(err)
	}
	roundTrip, _, err := FromCloudEvent(event, state)
	if err != nil || !roundTrip.Equal(message) {
		t.Fatalf("round trip = %#v, %v", roundTrip, err)
	}
}

type messageSpec struct {
	payload        []byte
	contentType    string
	recordedAt     time.Time
	metadata       map[string]string
	correlation    string
	causation      string
	tenant         string
	partition      string
	globalPosition golibeventsourcing.GlobalPosition
}

func newMessage(t *testing.T, spec messageSpec) golibeventsourcing.Message {
	t.Helper()
	if spec.contentType == "" {
		spec.contentType = "application/octet-stream"
	}
	if spec.recordedAt.IsZero() {
		spec.recordedAt = time.Date(2026, 8, 9, 1, 2, 3, 0, time.UTC)
	}
	stream, err := golibeventsourcing.NewStreamID("order", "A/123")
	if err != nil {
		t.Fatal(err)
	}
	encoded, err := golibeventsourcing.NewEncodedEvent(golibeventsourcing.EncodedEventInput{
		Name: "order.created", Version: 2, ContentType: spec.contentType, Payload: spec.payload,
	})
	if err != nil {
		t.Fatal(err)
	}
	pending, err := golibeventsourcing.NewPendingMessage(golibeventsourcing.PendingMessageInput{
		ID: "message-1", Stream: stream, Event: encoded, RecordedAt: spec.recordedAt, Metadata: spec.metadata,
		CorrelationID: spec.correlation, CausationID: spec.causation, Tenant: spec.tenant, Partition: spec.partition,
	})
	if err != nil {
		t.Fatal(err)
	}
	message, err := golibeventsourcing.NewMessage(golibeventsourcing.MessageInput{Pending: pending, StreamVersion: 4, GlobalPosition: spec.globalPosition})
	if err != nil {
		t.Fatal(err)
	}
	return message
}

func cloneEvent(t *testing.T, event cloudevents.Event, mutate func(*cloudevents.Attributes), data cloudevents.Data) cloudevents.Event {
	t.Helper()
	attributes := cloudevents.Attributes{
		ID: event.ID(), Source: event.Source(), Type: event.Type(), Extensions: event.Extensions(),
	}
	attributes.DataContentType, _ = event.DataContentType()
	attributes.DataSchema, _ = event.DataSchema()
	attributes.Subject, _ = event.Subject()
	if occurredAt, present := event.Time(); present {
		attributes.Time = &occurredAt
	}
	if mutate != nil {
		mutate(&attributes)
	}
	cloned, err := cloudevents.NewEvent(attributes, data)
	if err != nil {
		t.Fatal(err)
	}
	return cloned
}
