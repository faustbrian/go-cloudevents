package golib_test

import (
	"github.com/faustbrian/go-cloudevents"
	golib "github.com/faustbrian/go-cloudevents/adapters/golib/v4"
	jsonschema "github.com/faustbrian/go-json-schema/v2"
	registry "github.com/faustbrian/go-schema-registry/v3"
	outbox "github.com/faustbrian/go-transactional-outbox/v2"
)

// These assignments guard the public nominal producer cohort, not aliases of
// legacy types that merely share the same fields.
var _ func(outbox.Envelope, golib.OutboxOptions) (cloudevents.Event, outbox.Envelope, golib.Report, error) = golib.OutboxToCloudEvent
var _ func(cloudevents.Event, outbox.Envelope) (outbox.Envelope, golib.Report, error) = golib.CloudEventToOutbox
var _ *jsonschema.Schema = golib.JSONSchemaValidator{}.Schema
var _ map[string]registry.Lookup = golib.RegistryJSONSchemaConfig{}.SchemaLookups
