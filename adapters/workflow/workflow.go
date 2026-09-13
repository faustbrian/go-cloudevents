// Package workflow maps durable Golib workflow history decisions to and from
// CloudEvents while retaining workflow-owned state outside the event.
package workflow

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	golibworkflow "github.com/faustbrian/go-workflow"
)

// ErrInvalidInput reports workflow state or a CloudEvent that cannot satisfy
// the adapter contract. Its zero value is not meaningful; callers should
// compare returned errors with errors.Is.
var ErrInvalidInput = cloudevents.ErrInvalidAdapterInput

// Loss describes one CloudEvent field that workflow history cannot represent.
// The zero value describes no useful loss; values are owned by their Report.
type Loss = cloudevents.AdapterLoss

// Report lists deliberate information loss during a conversion. Its zero
// value means the conversion was lossless, and callers own the returned value.
type Report = cloudevents.AdapterReport

// Options supplies CloudEvent identity that workflow history does not own. The
// zero value is invalid because both StableID and Source are required.
type Options struct {
	// StableID is the required stable CloudEvent ID and is copied into State.
	StableID string
	// Source is the required CloudEvent source URI-reference. It is not retained
	// in State because workflow history cannot represent it on reversal.
	Source string
}

// State retains workflow-owned fields that are not portable CloudEvent
// attributes. Callers own State and must return the matching value unchanged;
// its zero value is invalid because StableID and Sequence are required.
type State struct {
	// StableID binds State to one CloudEvent. Empty is invalid, and a mismatch is
	// rejected rather than merged.
	StableID string
	// Sequence is the workflow history order. Zero is invalid.
	Sequence uint64
	// Definition is the workflow definition reference. Its zero value is passed
	// to workflow validation and is invalid for history kinds that require it.
	Definition golibworkflow.DefinitionReference
	// SuccessorID retains the workflow successor identifier. Empty means absent.
	SuccessorID string
	// StepName retains the workflow step name. Empty means absent.
	StepName string
	// Attempt retains the workflow attempt count. Zero means no attempt was set.
	Attempt uint32
	// IdempotencyKey retains the workflow idempotency key. Empty means absent.
	IdempotencyKey string
	// DueAt retains the workflow due time. The zero time means absent.
	DueAt time.Time
	// Code retains the workflow outcome or failure code. Empty means absent.
	Code string
	// Retryable retains whether workflow processing may be retried. False is the
	// zero value and means retryability was not asserted.
	Retryable bool
	// DataWasNil distinguishes nil payload data from a non-nil empty payload.
	// False is the zero value and means the source payload was non-nil.
	DataWasNil bool
}

// ToCloudEvent converts history into a CloudEvent and returns caller-owned
// State containing workflow-only fields. Payload bytes and occurrence time are
// copied, no input storage is retained, and invalid identity fails conversion.
func ToCloudEvent(history golibworkflow.HistoryEvent, options Options) (cloudevents.Event, State, Report, error) {
	if options.StableID == "" || options.Source == "" || history.Sequence() == 0 {
		return cloudevents.Event{}, State{}, Report{}, fmt.Errorf("%w: workflow mapping", ErrInvalidInput)
	}
	occurredAt := history.OccurredAt()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{ID: options.StableID, Source: options.Source, Type: eventType(history.Kind()), Subject: history.InstanceID(), Time: &occurredAt}, cloudevents.NewBinaryData(history.Data()))
	if err != nil {
		return cloudevents.Event{}, State{}, Report{}, err
	}
	return event, State{StableID: options.StableID, Sequence: history.Sequence(), Definition: history.Definition(), SuccessorID: history.SuccessorID(), StepName: history.StepName(), Attempt: history.Attempt(), IdempotencyKey: history.IdempotencyKey(), DueAt: history.DueAt(), Code: history.Code(), Retryable: history.Retryable(), DataWasNil: history.Data() == nil}, Report{}, nil
}

// FromCloudEvent reconstructs workflow history from portable event fields and
// caller-retained State. A missing or mismatched State identity fails instead
// of merging histories; Report names event fields the workflow model discards.
func FromCloudEvent(event cloudevents.Event, state State) (golibworkflow.HistoryEvent, Report, error) {
	if state.StableID == "" || state.Sequence == 0 || event.ID() != state.StableID {
		return golibworkflow.HistoryEvent{}, Report{}, fmt.Errorf("%w: workflow target", ErrInvalidInput)
	}
	kind, err := parseEventType(event.Type())
	if err != nil {
		return golibworkflow.HistoryEvent{}, Report{}, err
	}
	instanceID, present := event.Subject()
	occurredAt, hasTime := event.Time()
	if !present || !hasTime || !event.Data().Present() {
		return golibworkflow.HistoryEvent{}, Report{}, fmt.Errorf("%w: workflow portable fields", ErrInvalidInput)
	}
	history, err := golibworkflow.NewHistoryEvent(golibworkflow.HistoryEventSpec{Sequence: state.Sequence, InstanceID: instanceID, Kind: kind, OccurredAt: occurredAt, Definition: state.Definition, SuccessorID: state.SuccessorID, StepName: state.StepName, Attempt: state.Attempt, IdempotencyKey: state.IdempotencyKey, DueAt: state.DueAt, Code: state.Code, Retryable: state.Retryable, Data: adapter.RestoreRetainedNil(event.Data(), state.DataWasNil)})
	if err != nil {
		return golibworkflow.HistoryEvent{}, Report{}, err
	}
	report := Report{Losses: []Loss{{Field: "source", Reason: "not represented by workflow history"}}}
	if _, present := event.DataContentType(); present {
		report.Losses = append(report.Losses, Loss{Field: "datacontenttype", Reason: "not represented by workflow history"})
	}
	if _, present := event.DataSchema(); present {
		report.Losses = append(report.Losses, Loss{Field: "dataschema", Reason: "not represented by workflow history"})
	}
	adapter.AppendDataKindLoss(event, &report, "workflow history", false)
	adapter.AppendExtensionLosses(event, &report, "workflow history", nil)
	return history, report, nil
}

func eventType(kind golibworkflow.EventKind) string {
	return "golib.workflow.history." + strconv.FormatUint(uint64(kind), 10)
}

func parseEventType(value string) (golibworkflow.EventKind, error) {
	const prefix = "golib.workflow.history."
	if !strings.HasPrefix(value, prefix) {
		return 0, fmt.Errorf("%w: workflow type", ErrInvalidInput)
	}
	parsed, err := strconv.ParseUint(strings.TrimPrefix(value, prefix), 10, 8)
	if err != nil {
		return 0, fmt.Errorf("%w: workflow type: %w", ErrInvalidInput, err)
	}
	kind := golibworkflow.EventKind(parsed)
	if value != eventType(kind) {
		return 0, fmt.Errorf("%w: non-canonical workflow type", ErrInvalidInput)
	}
	return kind, nil
}
