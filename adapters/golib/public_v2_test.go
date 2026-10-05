package golib_test

import (
	"bytes"
	"testing"
	"time"

	golib "github.com/faustbrian/go-cloudevents/adapters/golib/v3"
	workflow "github.com/faustbrian/go-workflow/v2"
)

func TestPublicWorkflowV2FacadeComposition(t *testing.T) {
	reference, err := workflow.NewDefinitionReference("orders", "v1", "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef")
	if err != nil {
		t.Fatal(err)
	}
	payload := []byte("order-created")
	occurredAt := time.Date(2026, 10, 4, 12, 0, 0, 0, time.UTC)
	history, err := workflow.NewHistoryEvent(workflow.HistoryEventSpec{
		Sequence: 1, InstanceID: "order-1", Kind: workflow.EventInstanceStarted,
		OccurredAt: occurredAt, Definition: reference, Data: payload,
	})
	if err != nil {
		t.Fatal(err)
	}
	event, state, report, err := golib.WorkflowToCloudEvent(history, golib.WorkflowOptions{StableID: "event-1", Source: "/orders"})
	if err != nil || len(report.Losses) != 0 {
		t.Fatalf("WorkflowToCloudEvent() = %#v, %v", report, err)
	}
	if event.ID() != "event-1" || state.StableID != event.ID() || state.Definition != reference {
		t.Fatal("retained workflow identity differs from source")
	}
	payload[0] = 'X'
	historyData := history.Data()
	historyData[0] = 'Y'
	if !bytes.Equal(event.Data().Bytes(), []byte("order-created")) {
		t.Fatal("event payload aliases source storage")
	}
	var restored workflow.HistoryEvent
	restored, report, err = golib.CloudEventToWorkflow(event, state)
	if err != nil || len(report.Losses) != 1 || report.Losses[0].Field != "source" {
		t.Fatalf("CloudEventToWorkflow() = %#v, %v", report, err)
	}
	if restored.Sequence() != history.Sequence() || restored.InstanceID() != history.InstanceID() ||
		restored.Kind() != history.Kind() || restored.Definition() != reference ||
		!restored.OccurredAt().Equal(occurredAt) || !bytes.Equal(restored.Data(), []byte("order-created")) {
		t.Fatal("public Workflow/v2 history did not round trip")
	}
	restoredData := restored.Data()
	restoredData[0] = 'Z'
	if !bytes.Equal(event.Data().Bytes(), []byte("order-created")) {
		t.Fatal("restored payload aliases event storage")
	}
}
