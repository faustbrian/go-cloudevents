package cloudevents

import (
	"errors"
	"strings"
)

// ErrInvalidEvent identifies a CloudEvent that violates the supported
// specification contract.
var ErrInvalidEvent = errors.New("cloudevents: invalid event")

// ErrInvalidData identifies event data that cannot be represented by its
// declared runtime kind.
var ErrInvalidData = errors.New("cloudevents: invalid data")

// ErrInvalidAttribute identifies a context value outside the CloudEvents type
// system.
var ErrInvalidAttribute = errors.New("cloudevents: invalid attribute")

// ErrLimitExceeded identifies an input rejected before an unchecked
// allocation or conversion.
var ErrLimitExceeded = errors.New("cloudevents: limit exceeded")

// ErrInvalidAdapterInput classifies an absent or malformed optional adapter
// input. Target adapters wrap this sentinel while preserving the safe cause.
var ErrInvalidAdapterInput = errors.New("cloudevents golib adapter: invalid input")

// ErrMetadataCollision reports canonical metadata that conflicts with an
// existing CloudEvents value and therefore cannot be overwritten.
var ErrMetadataCollision = errors.New("cloudevents golib adapter: metadata collision")

// ErrUntrustedMetadata reports an attempt to adopt inbound identity metadata
// without an explicit trust decision.
var ErrUntrustedMetadata = errors.New("cloudevents golib adapter: metadata is untrusted")

// ErrSchemaViolation reports a valid event payload that does not satisfy the
// explicitly selected schema.
var ErrSchemaViolation = errors.New("cloudevents golib adapter: schema violation")

// ErrSchemaMapping reports an absent, conflicting, or unsupported explicit
// schema selection in an optional adapter.
var ErrSchemaMapping = errors.New("cloudevents golib adapter: schema mapping")

// AdapterLoss describes one value that a selected adapter target cannot
// represent. It contains field names and reasons, never field values.
type AdapterLoss struct {
	Field  string
	Reason string
}

// AdapterReport makes every optional-adapter conversion loss explicit.
type AdapterReport struct {
	Losses []AdapterLoss
}

// IssueCode is a stable, value-free validation diagnostic.
type IssueCode string

const (
	IssueRequired            IssueCode = "required"
	IssueInvalidString       IssueCode = "invalid_string"
	IssueInvalidURIReference IssueCode = "invalid_uri_reference"
	IssueAbsoluteURIRequired IssueCode = "absolute_uri_required"
	IssueInvalidMediaType    IssueCode = "invalid_media_type"
	IssueInvalidName         IssueCode = "invalid_name"
	IssueReservedName        IssueCode = "reserved_name"
	IssueInvalidAttribute    IssueCode = "invalid_attribute"
)

// Issue identifies one invalid field without retaining or disclosing its
// value.
type Issue struct {
	Field string
	Code  IssueCode
}

// ValidationError contains canonical, field-sorted validation diagnostics.
type ValidationError struct {
	issues []Issue
}

// Error returns a deterministic diagnostic that never contains rejected
// values.
func (e *ValidationError) Error() string {
	var message strings.Builder
	message.WriteString(ErrInvalidEvent.Error())
	message.WriteString(": ")
	for index, issue := range e.issues {
		if index > 0 {
			message.WriteString("; ")
		}
		message.WriteString(issue.Field)
		message.WriteByte(' ')
		message.WriteString(string(issue.Code))
	}
	return message.String()
}

// Unwrap makes every ValidationError match ErrInvalidEvent.
func (e *ValidationError) Unwrap() error { return ErrInvalidEvent }

// Issues returns an owned copy of the canonical diagnostics.
func (e *ValidationError) Issues() []Issue {
	return append([]Issue(nil), e.issues...)
}
