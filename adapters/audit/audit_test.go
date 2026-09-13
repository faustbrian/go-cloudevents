package audit_test

import (
	"errors"
	"testing"
	"time"

	"github.com/faustbrian/go-audit"
	"github.com/faustbrian/go-cloudevents"
	cloudaudit "github.com/faustbrian/go-cloudevents/adapters/audit"
)

func TestAuditMetadataRoundTripSelectsSafeFieldsAndReportsLoss(t *testing.T) {
	t.Parallel()

	record := auditRecord(t, "tenant-a", "correlation-1", "cause-1")
	event, report, err := cloudaudit.AddMetadata(baseEvent(t), record)
	if err != nil {
		t.Fatal(err)
	}
	wantLosses := []string{
		"audit.occurred_at", "audit.recorded_at", "audit.reason_code", "audit.description",
		"audit.actor", "audit.subject", "audit.context.request_id", "audit.context.trace_id",
		"audit.context.idempotency_id", "audit.context.source_service", "audit.context.source_version",
		"audit.context.environment", "audit.context.network_origin", "audit.context.user_agent",
		"audit.changes", "audit.policy", "audit.integrity", "audit.attributes", "audit.redaction_applied",
	}
	if len(report.Losses) != len(wantLosses) {
		t.Fatalf("AddMetadata() losses = %#v; want fields %#v", report.Losses, wantLosses)
	}
	for index, field := range wantLosses {
		if report.Losses[index].Field != field || report.Losses[index].Reason == "" {
			t.Fatalf("loss %d = %#v; want field %q with reason", index, report.Losses[index], field)
		}
	}
	metadata, err := cloudaudit.ExtractMetadata(event, true)
	if err != nil || metadata.RecordID != record.ID() || metadata.Action != record.Action() ||
		metadata.Outcome != record.Outcome() || metadata.Tenant.Value() != "tenant-a" ||
		metadata.Correlation.CorrelationID.String() != "correlation-1" ||
		metadata.Correlation.CausationID.String() != "cause-1" {
		t.Fatalf("ExtractMetadata() = %#v, %v", metadata, err)
	}
}

func TestAuditMetadataRoundTripAllowsAbsentContextIdentifiers(t *testing.T) {
	t.Parallel()

	event, _, err := cloudaudit.AddMetadata(baseEvent(t), auditRecord(t, "", "", ""))
	if err != nil {
		t.Fatal(err)
	}
	metadata, err := cloudaudit.ExtractMetadata(event, true)
	if err != nil || metadata.Tenant.Valid() || metadata.Correlation.CorrelationID != "" || metadata.Correlation.CausationID != "" {
		t.Fatalf("ExtractMetadata() = %#v, %v; want empty context", metadata, err)
	}
}

func TestAuditAddRejectsInvalidRecordTenantAndCollision(t *testing.T) {
	t.Parallel()

	if _, _, err := cloudaudit.AddMetadata(baseEvent(t), audit.Record{}); !errors.Is(err, cloudaudit.ErrInvalidInput) {
		t.Fatalf("zero record error = %v", err)
	}
	if _, _, err := cloudaudit.AddMetadata(baseEvent(t), auditRecord(t, "bad?tenant", "", "")); !errors.Is(err, cloudaudit.ErrInvalidInput) {
		t.Fatalf("invalid tenant error = %v", err)
	}
	other, err := cloudevents.NewStringAttribute("other")
	if err != nil {
		t.Fatal(err)
	}
	if _, _, err := cloudaudit.AddMetadata(
		eventWithExtensions(t, map[string]cloudevents.Attribute{"auditid": other}),
		auditRecord(t, "", "", ""),
	); !errors.Is(err, cloudaudit.ErrMetadataCollision) {
		t.Fatalf("collision error = %v", err)
	}
}

func TestAuditExtractRejectsUntrustedIncompleteAndMalformedMetadata(t *testing.T) {
	t.Parallel()

	if metadata, err := cloudaudit.ExtractMetadata(baseEvent(t), false); err != nil || metadata != (cloudaudit.Metadata{}) {
		t.Fatalf("absent ExtractMetadata() = %#v, %v", metadata, err)
	}
	complete := auditExtensions(t, "1")
	if _, err := cloudaudit.ExtractMetadata(eventWithExtensions(t, complete), false); !errors.Is(err, cloudaudit.ErrUntrustedMetadata) {
		t.Fatalf("untrusted metadata error = %v", err)
	}

	for _, name := range []string{"auditid", "auditaction", "auditoutcome"} {
		name := name
		t.Run(name+" wrong attribute kind", func(t *testing.T) {
			t.Parallel()
			if _, err := cloudaudit.ExtractMetadata(eventWithExtensions(t, map[string]cloudevents.Attribute{
				name: cloudevents.NewBooleanAttribute(true),
			}), true); !errors.Is(err, cloudaudit.ErrInvalidInput) {
				t.Fatalf("ExtractMetadata() error = %v", err)
			}
		})
	}

	for _, missing := range []string{"auditid", "auditaction", "auditoutcome"} {
		extensions := auditExtensions(t, "1")
		delete(extensions, missing)
		if _, err := cloudaudit.ExtractMetadata(eventWithExtensions(t, extensions), true); !errors.Is(err, cloudaudit.ErrInvalidInput) {
			t.Fatalf("missing %s error = %v", missing, err)
		}
	}
	for _, outcome := range []string{"bad", "0", "5"} {
		if _, err := cloudaudit.ExtractMetadata(eventWithExtensions(t, auditExtensions(t, outcome)), true); !errors.Is(err, cloudaudit.ErrInvalidInput) {
			t.Fatalf("outcome %q error = %v", outcome, err)
		}
	}

	badCorrelation, _ := cloudevents.NewStringAttribute("bad.value")
	withBadCorrelation := auditExtensions(t, "1")
	withBadCorrelation["correlationid"] = badCorrelation
	if _, err := cloudaudit.ExtractMetadata(eventWithExtensions(t, withBadCorrelation), true); !errors.Is(err, cloudaudit.ErrInvalidInput) {
		t.Fatalf("bad correlation error = %v", err)
	}
	badTenant, _ := cloudevents.NewStringAttribute("bad?tenant")
	withBadTenant := auditExtensions(t, "1")
	withBadTenant["tenantid"] = badTenant
	if _, err := cloudaudit.ExtractMetadata(eventWithExtensions(t, withBadTenant), true); !errors.Is(err, cloudaudit.ErrInvalidInput) {
		t.Fatalf("bad tenant error = %v", err)
	}
}

func auditRecord(t *testing.T, tenant, correlationID, causationID string) audit.Record {
	t.Helper()
	now := time.Date(2026, 9, 9, 12, 0, 0, 0, time.UTC)
	builder, err := audit.NewBuilder(audit.BuilderConfig{
		Clock: func() time.Time { return now }, IDGenerator: func() (string, error) { return "audit-1", nil },
	})
	if err != nil {
		t.Fatal(err)
	}
	record, err := builder.Build(audit.RecordInput{
		OccurredAt: now, Action: "order.create", Outcome: audit.OutcomeSucceeded,
		Actor:   audit.ActorInput{Kind: audit.ActorService, ID: "orders"},
		Subject: audit.SubjectInput{Type: "order", ID: "A-123"},
		Context: audit.ContextInput{TenantID: tenant, CorrelationID: correlationID, CausationID: causationID},
		Changes: audit.ChangeSetInput{NoChange: true}, Policy: audit.PolicyMetadata{PolicyID: "audit", Version: "v1"},
	})
	if err != nil {
		t.Fatal(err)
	}
	return record
}

func auditExtensions(t *testing.T, outcome string) map[string]cloudevents.Attribute {
	t.Helper()
	return map[string]cloudevents.Attribute{
		"auditid":      stringAttribute(t, "audit-1"),
		"auditaction":  stringAttribute(t, "order.create"),
		"auditoutcome": stringAttribute(t, outcome),
	}
}

func stringAttribute(t *testing.T, value string) cloudevents.Attribute {
	t.Helper()
	attribute, err := cloudevents.NewStringAttribute(value)
	if err != nil {
		t.Fatal(err)
	}
	return attribute
}

func baseEvent(t *testing.T) cloudevents.Event {
	t.Helper()
	return eventWithExtensions(t, nil)
}

func eventWithExtensions(t *testing.T, extensions map[string]cloudevents.Attribute) cloudevents.Event {
	t.Helper()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "event-1", Source: "/source", Type: "example.created", Extensions: extensions,
	}, cloudevents.NewBinaryData([]byte("body")))
	if err != nil {
		t.Fatal(err)
	}
	return event
}
