package golib

import (
	"github.com/faustbrian/go-cloudevents"
	cloudeventsourcing "github.com/faustbrian/go-cloudevents/adapters/event-sourcing/v2"
	cloudoutbox "github.com/faustbrian/go-cloudevents/adapters/outbox/v2"
	cloudqueue "github.com/faustbrian/go-cloudevents/adapters/queue/v2"
	cloudworkflow "github.com/faustbrian/go-cloudevents/adapters/workflow/v2"
	eventsourcing "github.com/faustbrian/go-event-sourcing/v2"
	"github.com/faustbrian/go-queue/job"
	outbox "github.com/faustbrian/go-transactional-outbox/v2"
	"github.com/faustbrian/go-workflow/v2"
)

// EventSourcingOptions is the event-sourcing/v2 occurrence/source configuration alias.
type EventSourcingOptions = cloudeventsourcing.Options

// EventSourcingState is the event-sourcing/v2 retained store-state alias; callers must pair it with the event.
type EventSourcingState = cloudeventsourcing.State

// OutboxOptions is the outbox/v2 portable-context alias, not relay configuration.
type OutboxOptions = cloudoutbox.Options

// QueueOptions is the queue/v2 portable-context alias, not execution/settlement policy.
type QueueOptions = cloudqueue.Options

// WorkflowOptions is the workflow/v2 stable-ID/source alias, not a workflow start request.
type WorkflowOptions = cloudworkflow.Options

// WorkflowState is the workflow/v2 retained-history alias; its identity must match reversal.
type WorkflowState = cloudworkflow.State

// EventSourcingToCloudEvent delegates event-sourcing/v2 conversion and copied retained store state, without appending to a store.
func EventSourcingToCloudEvent(message eventsourcing.Message, options EventSourcingOptions) (cloudevents.Event, EventSourcingState, Report, error) {
	return cloudeventsourcing.ToCloudEvent(message, options)
}

// CloudEventToEventSourcing delegates event-sourcing/v2 exact retained-identity/trust checks and loss reporting, without persisting.
func CloudEventToEventSourcing(event cloudevents.Event, state EventSourcingState) (eventsourcing.Message, Report, error) {
	return cloudeventsourcing.FromCloudEvent(event, state)
}

// OutboxToCloudEvent delegates outbox/v2 payload/context mapping and retained-state copying, without publishing or starting a relay.
func OutboxToCloudEvent(envelope outbox.Envelope, options OutboxOptions) (cloudevents.Event, outbox.Envelope, Report, error) {
	return cloudoutbox.ToCloudEvent(envelope, options)
}

// CloudEventToOutbox delegates outbox/v2 identity checks/retained-state reversal and loss reporting, without a transaction.
func CloudEventToOutbox(event cloudevents.Event, state outbox.Envelope) (outbox.Envelope, Report, error) {
	return cloudoutbox.FromCloudEvent(event, state)
}

// QueueToCloudEvent delegates queue/v2 conversion and copied execution state, without enqueueing.
func QueueToCloudEvent(message job.Message, options QueueOptions) (cloudevents.Event, job.Message, Report, error) {
	return cloudqueue.ToCloudEvent(message, options)
}

// CloudEventToQueue delegates queue/v2 collision-checked reversal and loss reporting, without execution or settlement.
func CloudEventToQueue(event cloudevents.Event, state job.Message) (job.Message, Report, error) {
	return cloudqueue.FromCloudEvent(event, state)
}

// WorkflowToCloudEvent delegates workflow/v2 history conversion with copied occurrence/payload state, without starting a workflow.
func WorkflowToCloudEvent(history workflow.HistoryEvent, options WorkflowOptions) (cloudevents.Event, WorkflowState, Report, error) {
	return cloudworkflow.ToCloudEvent(history, options)
}

// CloudEventToWorkflow delegates workflow/v2 stable identity/sequence reconstruction and loss reporting, without appending history.
func CloudEventToWorkflow(event cloudevents.Event, state WorkflowState) (workflow.HistoryEvent, Report, error) {
	return cloudworkflow.FromCloudEvent(event, state)
}
