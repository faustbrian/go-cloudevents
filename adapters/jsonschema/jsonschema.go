// Package jsonschema adapts one caller-compiled Golib JSON Schema to the
// CloudEvents opt-in SchemaValidator boundary without registry lookup.
package jsonschema

import (
	"context"
	"fmt"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	golibjsonschema "github.com/faustbrian/go-json-schema"
)

var (
	// ErrSchemaViolation reports that a JSON payload does not satisfy the
	// configured schema. The error aliases the CloudEvents facade sentinel so
	// callers can use errors.Is across either package boundary.
	ErrSchemaViolation = cloudevents.ErrSchemaViolation
	// ErrSchemaMapping reports missing, mismatched, or unsupported schema
	// configuration. The error aliases the CloudEvents facade sentinel.
	ErrSchemaMapping = cloudevents.ErrSchemaMapping
)

// Validator validates JSON payloads against one caller-owned compiled schema.
// It performs no registry or network lookup and is safe for concurrent use
// when Schema is safe for concurrent validation.
type Validator struct {
	// URI is the exact trusted dataschema identifier accepted by Validate.
	// Empty values are invalid; no normalization or prefix matching occurs.
	URI string
	// Schema is the caller-compiled schema used for validation. The caller owns
	// its construction, lifecycle, and any trust decision for its contents.
	Schema *golibjsonschema.Schema
}

// Validate implements cloudevents.SchemaValidator. It requires a non-nil
// context, an exact URI match, and a JSON-compatible content type. Mapping and
// content-type failures wrap ErrSchemaMapping, while structurally valid JSON
// that violates Schema returns ErrSchemaViolation. Parser, limit, and context
// errors from Schema are returned unchanged.
func (validator Validator) Validate(ctx context.Context, uri, contentType string, data []byte) error {
	if ctx == nil {
		return cloudevents.ErrContextRequired
	}
	if validator.URI == "" || validator.Schema == nil || uri != validator.URI {
		return ErrSchemaMapping
	}
	if !adapter.IsJSONContentType(contentType) {
		return fmt.Errorf("%w: non-JSON content type", ErrSchemaMapping)
	}
	result, err := validator.Schema.Validate(ctx, data)
	if err != nil {
		return err
	}
	if !result.Valid {
		return ErrSchemaViolation
	}
	return nil
}
