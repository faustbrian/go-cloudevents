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

// ErrLimitExceeded identifies a selected byte/count/depth admission failure.
// It does not certify that no temporary parsing or serialization allocation occurred.
var ErrLimitExceeded = errors.New("cloudevents: limit exceeded")

// ErrInvalidAdapterInput classifies an absent or malformed optional adapter
// input. Target adapters may wrap a more specific cause; metadata names and
// collaborator diagnostics still require caller protection before exposure.
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
// represent. Package-generated reports contain field names/reasons, not values;
// manually constructed text and untrusted identifier names are not sanitized.
type AdapterLoss struct {
	// Field names the affected path; identifier names are not redacted.
	Field string
	// Reason describes the representation omission/change, not its value; custom text remains caller-owned.
	Reason string
}

// AdapterReport makes every optional-adapter conversion loss explicit.
type AdapterReport struct {
	// Losses contains caller-owned ordered diagnostics; an empty slice reports no observed conversion loss.
	Losses []AdapterLoss
}

// IssueCode is a stable, value-free validation diagnostic.
type IssueCode string

const (
	// IssueRequired marks missing required context.
	IssueRequired IssueCode = "required"
	// IssueInvalidString marks text outside the selected CloudEvents character policy.
	IssueInvalidString IssueCode = "invalid_string"
	// IssueInvalidURIReference marks an invalid serialized source URI-reference.
	IssueInvalidURIReference IssueCode = "invalid_uri_reference"
	// IssueAbsoluteURIRequired marks a schema URI that is not absolute serialized URI syntax.
	IssueAbsoluteURIRequired IssueCode = "absolute_uri_required"
	// IssueInvalidMediaType marks an unparseable media type.
	IssueInvalidMediaType IssueCode = "invalid_media_type"
	// IssueInvalidName marks an extension name outside lowercase ASCII alphanumeric syntax.
	IssueInvalidName IssueCode = "invalid_name"
	// IssueReservedName marks a context/data name used as an extension.
	IssueReservedName IssueCode = "reserved_name"
	// IssueInvalidAttribute marks an unassigned or invalid abstract attribute/selected extension.
	IssueInvalidAttribute IssueCode = "invalid_attribute"
)

// Issue identifies one invalid field without retaining or disclosing its
// value.
type Issue struct {
	// Field names the affected path; identifier names are not redacted.
	Field string
	// Code is the failed validation category, not a rejected value.
	Code IssueCode
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
