package workflow

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/faustbrian/go-cloudevents"
	golibworkflow "github.com/faustbrian/go-workflow"
)

func TestWorkflowRoundTripPreservesOwnedStateWithoutAliasing(t *testing.T) {
	t.Parallel()

	payload := []byte("body")
	history := workflowHistory(t, payload)
	event, state, report, err := ToCloudEvent(history, Options{StableID: "event-1", Source: "/workflow"})
	if err != nil || len(report.Losses) != 0 {
		t.Fatalf("ToCloudEvent() = %#v, %v", report, err)
	}
	payload[0] = 'X'
	if string(event.Data().Bytes()) != "body" {
		t.Fatal("event payload aliases source payload")
	}
	roundTrip, reverseReport, err := FromCloudEvent(event, state)
	if err != nil || len(reverseReport.Losses) != 1 || roundTrip.Sequence() != history.Sequence() ||
		roundTrip.InstanceID() != history.InstanceID() || !bytes.Equal(roundTrip.Data(), []byte("body")) {
		t.Fatalf("FromCloudEvent() = %#v, %#v, %v", roundTrip, reverseReport, err)
	}
}

func TestWorkflowRejectsInvalidBoundaries(t *testing.T) {
	t.Parallel()

	history := workflowHistory(t, []byte("body"))
	for name, input := range map[string]struct {
		history golibworkflow.HistoryEvent
		options Options
	}{
		"missing ID":     {history: history, options: Options{Source: "/source"}},
		"missing source": {history: history, options: Options{StableID: "id"}},
		"zero history":   {options: Options{StableID: "id", Source: "/source"}},
	} {
		if _, _, _, err := ToCloudEvent(input.history, input.options); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("%s error = %v", name, err)
		}
	}
	if _, _, _, err := ToCloudEvent(history, Options{StableID: "id", Source: "bad\nsource"}); err == nil {
		t.Fatal("invalid source error = nil")
	}
	event, state, _, err := ToCloudEvent(history, Options{StableID: "id", Source: "/source"})
	if err != nil {
		t.Fatal(err)
	}
	for name, changed := range map[string]State{
		"empty":            {},
		"missing ID":       func() State { value := state; value.StableID = ""; return value }(),
		"missing sequence": func() State { value := state; value.Sequence = 0; return value }(),
		"different ID":     func() State { value := state; value.StableID = "other"; return value }(),
	} {
		if _, _, err := FromCloudEvent(event, changed); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("%s state error = %v", name, err)
		}
	}
	if _, _, err := FromCloudEvent(cloudevents.Event{}, state); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("zero event error = %v", err)
	}
	occurredAt, _ := event.Time()
	for name, candidate := range map[string]cloudevents.Event{
		"wrong type":      workflowEvent(t, event.ID(), "other", eventSubject(event), &occurredAt, event.Data()),
		"missing subject": workflowEvent(t, event.ID(), event.Type(), "", &occurredAt, event.Data()),
		"missing time":    workflowEvent(t, event.ID(), event.Type(), eventSubject(event), nil, event.Data()),
		"missing data":    workflowEvent(t, event.ID(), event.Type(), eventSubject(event), &occurredAt, cloudevents.Data{}),
	} {
		if _, _, err := FromCloudEvent(candidate, state); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("%s error = %v", name, err)
		}
	}
	for _, value := range []string{"golib.workflow.history.0", "golib.workflow.history.01", "golib.workflow.history.+1", "golib.workflow.history.999"} {
		candidate := workflowEvent(t, event.ID(), value, eventSubject(event), &occurredAt, event.Data())
		if _, _, err := FromCloudEvent(candidate, state); err == nil {
			t.Fatalf("invalid type %q error = nil", value)
		}
	}
}

func TestWorkflowPreservesNilAndEmptyPayloads(t *testing.T) {
	t.Parallel()

	for name, payload := range map[string][]byte{"nil": nil, "empty": {}} {
		payload := payload
		t.Run(name, func(t *testing.T) {
			event, state, _, err := ToCloudEvent(workflowHistory(t, payload), Options{StableID: "id", Source: "/workflow"})
			if err != nil {
				t.Fatal(err)
			}
			roundTrip, _, err := FromCloudEvent(event, state)
			if err != nil || (payload == nil) != (roundTrip.Data() == nil) || len(roundTrip.Data()) != 0 {
				t.Fatalf("round-trip payload = %#v, %v", roundTrip.Data(), err)
			}
		})
	}
}

func TestWorkflowReportsEveryDiscardedCloudEventField(t *testing.T) {
	t.Parallel()

	_, state, _, err := ToCloudEvent(workflowHistory(t, []byte("body")), Options{StableID: "id", Source: "/workflow"})
	if err != nil {
		t.Fatal(err)
	}
	data, err := cloudevents.NewJSONData([]byte(`{"ok":true}`))
	if err != nil {
		t.Fatal(err)
	}
	extra, _ := cloudevents.NewStringAttribute("value")
	occurredAt := time.Now()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "id", Source: "/source", Type: eventType(golibworkflow.EventInstanceStarted),
		Subject: "workflow-1", Time: &occurredAt, DataContentType: "application/json",
		DataSchema: "https://example.test/schema", Extensions: map[string]cloudevents.Attribute{"extra": extra},
	}, data)
	if err != nil {
		t.Fatal(err)
	}
	if _, report, err := FromCloudEvent(event, state); err != nil || len(report.Losses) != 5 {
		t.Fatalf("loss report = %#v, %v", report, err)
	}
}

func workflowHistory(t *testing.T, payload []byte) golibworkflow.HistoryEvent {
	t.Helper()
	reference, err := golibworkflow.NewDefinitionReference("orders", "v1", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	history, err := golibworkflow.NewHistoryEvent(golibworkflow.HistoryEventSpec{
		Sequence: 1, InstanceID: "workflow-1", Kind: golibworkflow.EventInstanceStarted,
		OccurredAt: time.Date(2026, 8, 10, 1, 2, 3, 0, time.UTC), Definition: reference, Data: payload,
	})
	if err != nil {
		t.Fatal(err)
	}
	return history
}

func workflowEvent(t *testing.T, id, eventTypeValue, subject string, occurredAt *time.Time, data cloudevents.Data) cloudevents.Event {
	t.Helper()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: id, Source: "/source", Type: eventTypeValue, Subject: subject, Time: occurredAt,
	}, data)
	if err != nil {
		t.Fatal(err)
	}
	return event
}

func eventSubject(event cloudevents.Event) string {
	value, _ := event.Subject()
	return value
}
