package golib

import (
	"github.com/faustbrian/go-audit"
	"github.com/faustbrian/go-cloudevents"
	cloudaudit "github.com/faustbrian/go-cloudevents/adapters/audit/v2"
)

// AuditMetadata is the audit/v2 selected-subset alias, not a reconstructed canonical audit record.
type AuditMetadata = cloudaudit.Metadata

// AddAuditMetadata delegates audit/v2 trusted subset mapping and omission reporting; input records are not persisted/redacted here.
func AddAuditMetadata(event cloudevents.Event, record audit.Record) (cloudevents.Event, Report, error) {
	return cloudaudit.AddMetadata(event, record)
}

// ExtractAuditMetadata delegates audit/v2 explicit trust and completeness checks without reconstructing an audit record.
func ExtractAuditMetadata(event cloudevents.Event, trusted bool) (AuditMetadata, error) {
	return cloudaudit.ExtractMetadata(event, trusted)
}
