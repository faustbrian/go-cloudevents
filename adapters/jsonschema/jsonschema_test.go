package jsonschema_test

import (
	"context"
	"errors"
	"testing"

	"github.com/faustbrian/go-cloudevents"
	cloudjsonschema "github.com/faustbrian/go-cloudevents/adapters/jsonschema/v2"
	golibjsonschema "github.com/faustbrian/go-json-schema/v2"
)

func TestValidatorEnforcesExplicitSchemaMapping(t *testing.T) {
	schema := compileSchema(t, `{"type":"object","required":["id"]}`)
	const uri = "https://schemas.example/orders/v1"

	tests := map[string]cloudjsonschema.Validator{
		"empty URI":      {Schema: schema},
		"missing schema": {URI: uri},
		"different URI":  {URI: uri, Schema: schema},
	}
	for name, validator := range tests {
		t.Run(name, func(t *testing.T) {
			gotURI := uri
			if name == "different URI" {
				gotURI = "https://schemas.example/orders/v2"
			}
			if err := validator.Validate(context.Background(), gotURI, "application/json", []byte(`{"id":1}`)); !errors.Is(err, cloudjsonschema.ErrSchemaMapping) {
				t.Fatalf("mapping error = %v, want %v", err, cloudjsonschema.ErrSchemaMapping)
			}
		})
	}
}

func TestValidatorRequiresContextAndJSONContent(t *testing.T) {
	validator := cloudjsonschema.Validator{
		URI:    "schema",
		Schema: compileSchema(t, `{}`),
	}

	//lint:ignore SA1012 Nil is the contract under test.
	if err := validator.Validate(nil, "schema", "application/json", []byte(`{}`)); !errors.Is(err, cloudevents.ErrContextRequired) { //nolint:staticcheck // Nil is the contract under test.
		t.Fatalf("nil context error = %v, want %v", err, cloudevents.ErrContextRequired)
	}
	if err := validator.Validate(context.Background(), "schema", "text/plain", []byte(`{}`)); !errors.Is(err, cloudjsonschema.ErrSchemaMapping) {
		t.Fatalf("content type error = %v, want %v", err, cloudjsonschema.ErrSchemaMapping)
	}
}

func TestValidatorReportsPayloadOutcomes(t *testing.T) {
	validator := cloudjsonschema.Validator{
		URI:    "schema",
		Schema: compileSchema(t, `{"type":"object","required":["id"]}`),
	}

	if err := validator.Validate(context.Background(), "schema", "application/problem+json", []byte(`{"id":1}`)); err != nil {
		t.Fatalf("valid payload: %v", err)
	}
	if err := validator.Validate(context.Background(), "schema", "application/json", []byte(`{}`)); !errors.Is(err, cloudjsonschema.ErrSchemaViolation) {
		t.Fatalf("schema violation error = %v, want %v", err, cloudjsonschema.ErrSchemaViolation)
	}
	if err := validator.Validate(context.Background(), "schema", "application/json", []byte(`{"id":`)); err == nil || errors.Is(err, cloudjsonschema.ErrSchemaViolation) || errors.Is(err, cloudjsonschema.ErrSchemaMapping) {
		t.Fatalf("malformed JSON error = %v, want underlying parser error", err)
	}
}

func compileSchema(t *testing.T, definition string) *golibjsonschema.Schema {
	t.Helper()
	compiler, err := golibjsonschema.NewCompiler()
	if err != nil {
		t.Fatal(err)
	}
	schema, err := compiler.Compile(context.Background(), []byte(definition))
	if err != nil {
		t.Fatal(err)
	}
	return schema
}

func TestValidatorPreservesV2UnicodeAndCancellationBoundaries(t *testing.T) {
	validator := cloudjsonschema.Validator{
		URI: "schema", Schema: compileSchema(t, `{"type":"string","maxLength":2}`),
	}
	if err := validator.Validate(t.Context(), "schema", "application/json", []byte(`"é🙂"`)); err != nil {
		t.Fatalf("two Unicode code points rejected: %v", err)
	}
	if err := validator.Validate(t.Context(), "schema", "application/json", []byte(`"é🙂x"`)); !errors.Is(err, cloudjsonschema.ErrSchemaViolation) {
		t.Fatalf("three Unicode code points error = %v, want schema violation", err)
	}
	err := validator.Validate(t.Context(), "schema", "application/json", []byte{'"', 0xff, '"'})
	if !errors.Is(err, golibjsonschema.ErrInvalidJSON) {
		t.Fatalf("invalid UTF-8 error = %v, want underlying parser failure", err)
	}
	// The producer uses encoding/json string decoding: escaped unpaired
	// surrogates become a replacement rune, unlike invalid raw UTF-8 bytes.
	replacement := cloudjsonschema.Validator{
		URI: "schema", Schema: compileSchema(t, `{"enum":["�"]}`),
	}
	if err := replacement.Validate(t.Context(), "schema", "application/json", []byte(`"\ud800"`)); err != nil {
		t.Fatalf("producer replacement-rune semantics changed: %v", err)
	}
	ctx, cancel := context.WithCancel(t.Context())
	cancel()
	if err := validator.Validate(ctx, "schema", "application/json", []byte(`"é🙂"`)); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation error = %v, want context cancellation", err)
	}
}
