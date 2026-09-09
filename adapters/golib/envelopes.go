package golib

import (
	"github.com/faustbrian/go-cloudevents"
	cloudeventsourcing "github.com/faustbrian/go-cloudevents/adapters/event-sourcing"
	cloudoutbox "github.com/faustbrian/go-cloudevents/adapters/outbox"
	cloudqueue "github.com/faustbrian/go-cloudevents/adapters/queue"
	cloudworkflow "github.com/faustbrian/go-cloudevents/adapters/workflow"
	eventsourcing "github.com/faustbrian/go-event-sourcing"
	"github.com/faustbrian/go-queue/job"
	outbox "github.com/faustbrian/go-transactional-outbox"
	"github.com/faustbrian/go-workflow"
)

type EventSourcingOptions = cloudeventsourcing.Options
type EventSourcingState = cloudeventsourcing.State
type OutboxOptions = cloudoutbox.Options
type QueueOptions = cloudqueue.Options
type WorkflowOptions = cloudworkflow.Options
type WorkflowState = cloudworkflow.State

func EventSourcingToCloudEvent(message eventsourcing.Message, options EventSourcingOptions) (cloudevents.Event, EventSourcingState, Report, error) {
	return cloudeventsourcing.ToCloudEvent(message, options)
}

func CloudEventToEventSourcing(event cloudevents.Event, state EventSourcingState) (eventsourcing.Message, Report, error) {
	return cloudeventsourcing.FromCloudEvent(event, state)
}

func OutboxToCloudEvent(envelope outbox.Envelope, options OutboxOptions) (cloudevents.Event, outbox.Envelope, Report, error) {
	return cloudoutbox.ToCloudEvent(envelope, options)
}

func CloudEventToOutbox(event cloudevents.Event, state outbox.Envelope) (outbox.Envelope, Report, error) {
	return cloudoutbox.FromCloudEvent(event, state)
}

func QueueToCloudEvent(message job.Message, options QueueOptions) (cloudevents.Event, job.Message, Report, error) {
	return cloudqueue.ToCloudEvent(message, options)
}

func CloudEventToQueue(event cloudevents.Event, state job.Message) (job.Message, Report, error) {
	return cloudqueue.FromCloudEvent(event, state)
}

func WorkflowToCloudEvent(history workflow.HistoryEvent, options WorkflowOptions) (cloudevents.Event, WorkflowState, Report, error) {
	return cloudworkflow.ToCloudEvent(history, options)
}

func CloudEventToWorkflow(event cloudevents.Event, state WorkflowState) (workflow.HistoryEvent, Report, error) {
	return cloudworkflow.FromCloudEvent(event, state)
}
