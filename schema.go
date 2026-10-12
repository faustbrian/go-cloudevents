package cloudevents

import (
	"context"
	"errors"
	"reflect"
)

var (
	// ErrSchemaRequired reports that explicit validation was requested for an
	// event without a dataschema URI.
	ErrSchemaRequired = errors.New("cloudevents: schema URI required")
	// ErrDataRequired reports that explicit validation was requested for an
	// event without data.
	ErrDataRequired = errors.New("cloudevents: event data required")
	// ErrSchemaValidatorRequired reports a nil opt-in validator.
	ErrSchemaValidatorRequired = errors.New("cloudevents: schema validator required")
	// ErrContextRequired reports a nil context passed to an explicit operation.
	ErrContextRequired = errors.New("cloudevents: context required")
)

// SchemaValidator is implemented by an explicit JSON Schema or schema-registry
// adapter. ValidateSchema supplies an owned data copy; direct callers choose
// ownership. Implementations must apply trusted schema
// selection, byte/work limits and cooperative context handling. Calls are direct:
// panic recovery, retries, post-call cancellation checks and shutdown are not added.
// Shared implementations must support their callers' concurrency. The core
// package never resolves schemas or performs I/O itself.
type SchemaValidator interface {
	// Validate checks the exact supplied URI/type/payload under caller-owned trust
	// policy, returning raw errors which callers must protect before exposure.
	Validate(ctx context.Context, uri string, contentType string, data []byte) error
}

// ValidateSchema invokes a caller-supplied validator for an event that declares
// both dataschema and data. Nil contexts/typed-nil validators are rejected, and
// cancellation is checked before calling the validator with a data copy. The
// callback can still block/panic or succeed after cancellation; its raw error
// and external resource policy remain caller-owned.
func ValidateSchema(ctx context.Context, event Event, validator SchemaValidator) error {
	if ctx == nil {
		return ErrContextRequired
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if nilSchemaValidator(validator) {
		return ErrSchemaValidatorRequired
	}
	if err := event.Validate(); err != nil {
		return err
	}
	schema, present := event.DataSchema()
	if !present {
		return ErrSchemaRequired
	}
	data := event.Data()
	if !data.Present() {
		return ErrDataRequired
	}
	contentType, _ := event.DataContentType()
	return validator.Validate(ctx, schema, contentType, data.Bytes())
}

func nilSchemaValidator(validator SchemaValidator) bool {
	if validator == nil {
		return true
	}
	value := reflect.ValueOf(validator)
	switch value.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return value.IsNil()
	default:
		return false
	}
}
