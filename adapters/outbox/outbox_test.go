package outbox

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/faustbrian/go-cloudevents"
	goliboutbox "github.com/faustbrian/go-transactional-outbox"
)

func TestRoundTripRetainsOwnedStateWithoutAliasing(t *testing.T) {
	t.Parallel()

	now := time.Date(2026, 8, 9, 2, 3, 4, 0, time.UTC)
	occurredAt := now.Add(-time.Minute)
	envelope := goliboutbox.Envelope{
		ID: "outbox-1", Topic: "orders", Payload: []byte(`{"ok":true}`), PayloadVersion: 3,
		Metadata: map[string]string{"owner": "application"}, OrderingKey: "tenant-a",
		IdempotencyKey: "operation-1", AvailableAt: now, CreatedAt: now,
	}
	want := append([]byte(nil), envelope.CanonicalJSON()...)
	event, state, report, err := ToCloudEvent(envelope, Options{
		Source: "/outbox/orders", Type: "order.created", DataContentType: "application/json",
		DataSchema: "https://schemas.example/order.json", Subject: "order/A-1", OccurredAt: &occurredAt,
	})
	if err != nil || len(report.Losses) != 0 {
		t.Fatalf("ToCloudEvent() = %#v, %v", report, err)
	}
	envelope.Payload[0] = 'x'
	envelope.Metadata["owner"] = "changed"
	if bytes.Equal(event.Data().Bytes(), envelope.Payload) || state.Metadata["owner"] != "application" {
		t.Fatal("conversion aliases caller-owned state")
	}

	roundTrip, reverseReport, err := FromCloudEvent(event, state)
	if err != nil || !bytes.Equal(roundTrip.CanonicalJSON(), want) {
		t.Fatalf("FromCloudEvent() = %#v, %#v, %v", roundTrip, reverseReport, err)
	}
	wantLosses := map[string]bool{
		"source": true, "type": true, "datacontenttype": true, "data.kind": true,
		"dataschema": true, "subject": true, "time": true,
	}
	for _, loss := range reverseReport.Losses {
		delete(wantLosses, loss.Field)
	}
	if len(wantLosses) != 0 {
		t.Fatalf("loss report omitted fields: %#v", wantLosses)
	}
	state.Payload[0] = 'x'
	state.Metadata["owner"] = "retained-state-mutated"
	if bytes.Equal(roundTrip.Payload, state.Payload) || roundTrip.Metadata["owner"] != "application" {
		t.Fatal("reconstructed envelope aliases caller-owned retained state")
	}
}

func TestConversionRejectsInvalidBoundariesAndIdentifierCollision(t *testing.T) {
	t.Parallel()

	envelope := goliboutbox.Envelope{ID: "outbox-1", Topic: "topic", Payload: []byte("body"), PayloadVersion: 1}
	for name, candidate := range map[string]Options{
		"missing identity": {Source: "/source", Type: "type"},
		"invalid payload":  {Source: "/source", Type: "type", DataContentType: ";"},
		"invalid event":    {Source: "bad\nsource", Type: "type"},
	} {
		input := envelope
		if name == "missing identity" {
			input.ID = ""
		}
		if _, _, _, err := ToCloudEvent(input, candidate); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) && name != "invalid event" {
			t.Fatalf("%s error = %v", name, err)
		} else if err == nil {
			t.Fatalf("%s error = nil", name)
		}
	}
	if _, _, err := FromCloudEvent(cloudevents.Event{}, envelope); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("invalid event error = %v", err)
	}
	event := mustEvent(t, "event-1", cloudevents.NewBinaryData([]byte("body")))
	if _, _, err := FromCloudEvent(event, envelope); !errors.Is(err, cloudevents.ErrMetadataCollision) {
		t.Fatalf("identifier collision error = %v", err)
	}
	state := envelope
	state.ID = ""
	if _, _, err := FromCloudEvent(mustEvent(t, "event-1", cloudevents.Data{}), state); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("absent data error = %v", err)
	}
	state.Topic = ""
	if _, _, err := FromCloudEvent(event, state); !errors.Is(err, cloudevents.ErrInvalidAdapterInput) {
		t.Fatalf("invalid retained state error = %v", err)
	}
}

func TestRoundTripPreservesNilAndNonNilEmptyPayloads(t *testing.T) {
	t.Parallel()

	for name, payload := range map[string][]byte{"nil": nil, "empty": {}} {
		payload := payload
		t.Run(name, func(t *testing.T) {
			envelope := goliboutbox.Envelope{ID: "outbox-1", Topic: "orders", Payload: payload, PayloadVersion: 1}
			event, state, _, err := ToCloudEvent(envelope, Options{Source: "/outbox", Type: "order.created"})
			if err != nil {
				t.Fatal(err)
			}
			roundTrip, _, err := FromCloudEvent(event, state)
			if err != nil {
				t.Fatal(err)
			}
			if (payload == nil) != (roundTrip.Payload == nil) || len(roundTrip.Payload) != 0 {
				t.Fatalf("payload = %#v", roundTrip.Payload)
			}
		})
	}
}

func mustEvent(t *testing.T, id string, data cloudevents.Data) cloudevents.Event {
	t.Helper()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{ID: id, Source: "/source", Type: "type"}, data)
	if err != nil {
		t.Fatal(err)
	}
	return event
}
