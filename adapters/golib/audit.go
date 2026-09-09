package golib

import (
	"github.com/faustbrian/go-audit"
	"github.com/faustbrian/go-cloudevents"
	cloudaudit "github.com/faustbrian/go-cloudevents/adapters/audit"
)

type AuditMetadata = cloudaudit.Metadata

func AddAuditMetadata(event cloudevents.Event, record audit.Record) (cloudevents.Event, Report, error) {
	return cloudaudit.AddMetadata(event, record)
}

func ExtractAuditMetadata(event cloudevents.Event, trusted bool) (AuditMetadata, error) {
	return cloudaudit.ExtractMetadata(event, trusted)
}
