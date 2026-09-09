// Package audit maps a deliberately small, trusted subset of Golib audit
// metadata through CloudEvents without treating the event as an audit record.
package audit

import (
	"fmt"
	"strconv"

	golibaudit "github.com/faustbrian/go-audit"
	"github.com/faustbrian/go-cloudevents"
	cloudcorrelation "github.com/faustbrian/go-cloudevents/adapters/correlation"
	cloudtenancy "github.com/faustbrian/go-cloudevents/adapters/tenancy"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	golibcorrelation "github.com/faustbrian/go-correlation"
	golibtenancy "github.com/faustbrian/go-tenancy"
)

const (
	auditIDExtension      = "auditid"
	auditActionExtension  = "auditaction"
	auditOutcomeExtension = "auditoutcome"
)

var (
	// ErrInvalidInput reports a malformed event, audit record, or selected
	// metadata value. It aliases the root adapter error for errors.Is checks.
	ErrInvalidInput = cloudevents.ErrInvalidAdapterInput
	// ErrMetadataCollision reports that AddMetadata would overwrite a
	// non-equivalent caller-owned extension. Existing metadata is preserved.
	ErrMetadataCollision = cloudevents.ErrMetadataCollision
	// ErrUntrustedMetadata reports that audit extensions were present but the
	// caller had not established their trustworthiness.
	ErrUntrustedMetadata = cloudevents.ErrUntrustedMetadata
)

// Loss is one audit-owned field deliberately omitted from a CloudEvent. Its
// zero value carries no field or reason and has no effect by itself.
type Loss = cloudevents.AdapterLoss

// Report lists the information AddMetadata deliberately does not flatten. Its
// zero value means no loss was reported; callers own any resulting policy or
// persistence decision.
type Report = cloudevents.AdapterReport

// Metadata is the trusted audit subset recoverable from CloudEvent extensions.
// The zero value represents an event with no audit metadata.
type Metadata struct {
	// RecordID identifies the canonical audit record; empty only when Metadata
	// is the zero value returned for an event without audit extensions.
	RecordID string
	// Action is the audit action copied into the event. The adapter does not
	// interpret or authorize it.
	Action string
	// Outcome is the validated canonical audit outcome. Its zero value is not
	// accepted for present audit metadata.
	Outcome golibaudit.Outcome
	// Correlation contains validated identifiers selected from the audit
	// context. Missing identifiers retain their individual zero values.
	Correlation golibcorrelation.Values
	// Tenant is the validated tenant identifier selected from the audit
	// context. Its zero value means no tenant extension was present.
	Tenant golibtenancy.TenantID
}

// AddMetadata returns a new event containing the safe audit subset and a
// deterministic report of audit-owned fields that were not flattened. The
// input event and record remain caller-owned and unmodified. Existing
// equivalent extensions are accepted; conflicting extensions fail with
// ErrMetadataCollision. The caller remains responsible for the canonical audit
// record, authorization, redaction, storage, and delivery.
func AddMetadata(event cloudevents.Event, record golibaudit.Record) (cloudevents.Event, Report, error) {
	if record.ID() == "" {
		return cloudevents.Event{}, Report{}, fmt.Errorf("%w: audit record", ErrInvalidInput)
	}
	contextValue := record.Context()
	additions := map[string]string{
		auditIDExtension: record.ID(), auditActionExtension: record.Action(),
		auditOutcomeExtension: strconv.FormatUint(uint64(record.Outcome()), 10),
	}
	if contextValue.CorrelationID() != "" {
		additions[adapter.CorrelationIDExtension] = contextValue.CorrelationID()
	}
	if contextValue.CausationID() != "" {
		additions[adapter.CausationIDExtension] = contextValue.CausationID()
	}
	if value := contextValue.TenantID(); value != "" {
		tenant, err := golibtenancy.ParseTenantID(value)
		if err != nil {
			return cloudevents.Event{}, Report{}, fmt.Errorf("%w: audit tenant: %w", ErrInvalidInput, err)
		}
		additions[adapter.TenantIDExtension] = tenant.Value()
	}
	converted, err := adapter.AddStringExtensions(event, additions)
	if err != nil {
		return cloudevents.Event{}, Report{}, err
	}
	return converted, Report{Losses: []Loss{
		{Field: "audit.occurred_at", Reason: "canonical audit ownership"},
		{Field: "audit.recorded_at", Reason: "canonical audit ownership"},
		{Field: "audit.reason_code", Reason: "not flattened into CloudEvents"},
		{Field: "audit.description", Reason: "not flattened into CloudEvents"},
		{Field: "audit.actor", Reason: "not flattened into CloudEvents"},
		{Field: "audit.subject", Reason: "not flattened into CloudEvents"},
		{Field: "audit.context.request_id", Reason: "not flattened into CloudEvents"},
		{Field: "audit.context.trace_id", Reason: "not flattened into CloudEvents"},
		{Field: "audit.context.idempotency_id", Reason: "not flattened into CloudEvents"},
		{Field: "audit.context.source_service", Reason: "not flattened into CloudEvents"},
		{Field: "audit.context.source_version", Reason: "not flattened into CloudEvents"},
		{Field: "audit.context.environment", Reason: "not flattened into CloudEvents"},
		{Field: "audit.context.network_origin", Reason: "not flattened into CloudEvents"},
		{Field: "audit.context.user_agent", Reason: "not flattened into CloudEvents"},
		{Field: "audit.changes", Reason: "not flattened into CloudEvents"},
		{Field: "audit.policy", Reason: "not flattened into CloudEvents"},
		{Field: "audit.integrity", Reason: "not flattened into CloudEvents"},
		{Field: "audit.attributes", Reason: "not flattened into CloudEvents"},
		{Field: "audit.redaction_applied", Reason: "not flattened into CloudEvents"},
	}}, nil
}

// ExtractMetadata validates and returns the selected audit extensions after
// the caller has established their trustworthiness. An event with no audit
// extensions returns zero Metadata without requiring trust. Partial metadata,
// malformed selected context, and untrusted present metadata are rejected;
// extraction does not reconstruct or persist a canonical audit record.
func ExtractMetadata(event cloudevents.Event, trusted bool) (Metadata, error) {
	recordID, hasRecord, err := adapter.StringExtension(event, auditIDExtension)
	if err != nil {
		return Metadata{}, err
	}
	action, hasAction, err := adapter.StringExtension(event, auditActionExtension)
	if err != nil {
		return Metadata{}, err
	}
	outcomeValue, hasOutcome, err := adapter.StringExtension(event, auditOutcomeExtension)
	if err != nil {
		return Metadata{}, err
	}
	if !hasRecord && !hasAction && !hasOutcome {
		return Metadata{}, nil
	}
	if !trusted {
		return Metadata{}, ErrUntrustedMetadata
	}
	if !hasRecord || !hasAction || !hasOutcome {
		return Metadata{}, fmt.Errorf("%w: incomplete audit metadata", ErrInvalidInput)
	}
	parsed, err := strconv.ParseUint(outcomeValue, 10, 8)
	if err != nil || golibaudit.Outcome(parsed) < golibaudit.OutcomeSucceeded || golibaudit.Outcome(parsed) > golibaudit.OutcomeUnknown {
		return Metadata{}, fmt.Errorf("%w: audit outcome", ErrInvalidInput)
	}
	correlationValues, err := cloudcorrelation.Extract(event, true, golibcorrelation.Policy{})
	if err != nil {
		return Metadata{}, err
	}
	var tenant golibtenancy.TenantID
	if _, present := event.Extension(adapter.TenantIDExtension); present {
		tenant, err = cloudtenancy.Extract(event, true)
		if err != nil {
			return Metadata{}, err
		}
	}
	return Metadata{RecordID: recordID, Action: action, Outcome: golibaudit.Outcome(parsed), Correlation: correlationValues, Tenant: tenant}, nil
}
