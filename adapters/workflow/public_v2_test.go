package workflow_test

import (
	"bytes"
	"testing"
	"time"

	cloudworkflow "github.com/faustbrian/go-cloudevents/adapters/workflow/v2"
	workflow "github.com/faustbrian/go-workflow/v2"
)

func TestPublicWorkflowV2Composition(t *testing.T) {
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
	event, state, report, err := cloudworkflow.ToCloudEvent(history, cloudworkflow.Options{StableID: "event-1", Source: "/orders"})
	if err != nil || len(report.Losses) != 0 {
		t.Fatalf("ToCloudEvent() = %#v, %v", report, err)
	}
	var retainedDefinition workflow.DefinitionReference
	retainedDefinition = state.Definition
	if retainedDefinition != reference {
		t.Fatal("retained definition differs from source")
	}
	var restored workflow.HistoryEvent
	restored, report, err = cloudworkflow.FromCloudEvent(event, state)
	if err != nil || len(report.Losses) != 1 || report.Losses[0].Field != "source" {
		t.Fatalf("FromCloudEvent() = %#v, %v", report, err)
	}
	if restored.Sequence() != history.Sequence() || restored.InstanceID() != history.InstanceID() ||
		restored.Kind() != history.Kind() || restored.Definition() != reference ||
		!restored.OccurredAt().Equal(occurredAt) || !bytes.Equal(restored.Data(), payload) {
		t.Fatal("public Workflow/v2 history did not round trip")
	}
}
