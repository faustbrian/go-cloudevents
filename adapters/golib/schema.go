package golib

import (
	cloudjsonschema "github.com/faustbrian/go-cloudevents/adapters/jsonschema"
	cloudregistry "github.com/faustbrian/go-cloudevents/adapters/schema-registry"
)

var (
	ErrSchemaViolation = cloudjsonschema.ErrSchemaViolation
	ErrSchemaMapping   = cloudjsonschema.ErrSchemaMapping
)

type JSONSchemaValidator = cloudjsonschema.Validator
type RegistryJSONSchemaConfig = cloudregistry.JSONSchemaConfig
type RegistryJSONSchemaValidator = cloudregistry.JSONSchemaValidator

func NewRegistryJSONSchemaValidator(config RegistryJSONSchemaConfig) (RegistryJSONSchemaValidator, error) {
	return cloudregistry.NewJSONSchemaValidator(config)
}
