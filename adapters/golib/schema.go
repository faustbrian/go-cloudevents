package golib

import (
	cloudjsonschema "github.com/faustbrian/go-cloudevents/adapters/jsonschema/v2"
	cloudregistry "github.com/faustbrian/go-cloudevents/adapters/schema-registry/v3"
)

var (
	// ErrSchemaViolation preserves the target/root sentinel identity for errors.Is.
	ErrSchemaViolation = cloudjsonschema.ErrSchemaViolation
	// ErrSchemaMapping preserves the target/root sentinel identity for errors.Is.
	ErrSchemaMapping = cloudjsonschema.ErrSchemaMapping
)

// JSONSchemaValidator is the jsonschema/v2 caller-compiled validator alias; it performs no registry lookup.
type JSONSchemaValidator = cloudjsonschema.Validator

// RegistryJSONSchemaConfig is the schema-registry/v3 explicit trust/deadline configuration alias.
type RegistryJSONSchemaConfig = cloudregistry.JSONSchemaConfig

// RegistryJSONSchemaValidator is the schema-registry/v3 borrowed-cache validator alias; validation may resolve mappings.
type RegistryJSONSchemaValidator = cloudregistry.JSONSchemaValidator

// NewRegistryJSONSchemaValidator delegates schema-registry/v3 allowlist snapshot/deadline validation; it borrows cache/adapter lifecycle.
func NewRegistryJSONSchemaValidator(config RegistryJSONSchemaConfig) (RegistryJSONSchemaValidator, error) {
	return cloudregistry.NewJSONSchemaValidator(config)
}
