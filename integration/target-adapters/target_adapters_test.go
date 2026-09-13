package targetadapters_test

import (
	"testing"

	"github.com/faustbrian/go-cloudevents"
	cloudaudit "github.com/faustbrian/go-cloudevents/adapters/audit"
	cloudcorrelation "github.com/faustbrian/go-cloudevents/adapters/correlation"
	cloudeventsourcing "github.com/faustbrian/go-cloudevents/adapters/event-sourcing"
	cloudjsonschema "github.com/faustbrian/go-cloudevents/adapters/jsonschema"
	cloudkafka "github.com/faustbrian/go-cloudevents/adapters/kafka"
	cloudoutbox "github.com/faustbrian/go-cloudevents/adapters/outbox"
	cloudqueue "github.com/faustbrian/go-cloudevents/adapters/queue"
	cloudrabbitstream "github.com/faustbrian/go-cloudevents/adapters/rabbitstream"
	cloudregistry "github.com/faustbrian/go-cloudevents/adapters/schema-registry"
	cloudtelemetry "github.com/faustbrian/go-cloudevents/adapters/telemetry"
	cloudtenancy "github.com/faustbrian/go-cloudevents/adapters/tenancy"
	cloudworkflow "github.com/faustbrian/go-cloudevents/adapters/workflow"
)

func TestTargetAdaptersExposeIndependentIntegrationBoundaries(t *testing.T) {
	t.Parallel()

	var _ cloudevents.SchemaValidator = cloudjsonschema.Validator{}
	var _ cloudevents.SchemaValidator = cloudregistry.JSONSchemaValidator{}

	if cloudaudit.ErrInvalidInput != cloudcorrelation.ErrInvalidInput ||
		cloudcorrelation.ErrInvalidInput != cloudevents.ErrInvalidAdapterInput ||
		cloudtenancy.ErrUntrustedMetadata != cloudevents.ErrUntrustedMetadata ||
		cloudkafka.ErrMetadataCollision != cloudevents.ErrMetadataCollision {
		t.Fatal("target adapters do not share the canonical CloudEvents adapter errors")
	}

	_ = cloudeventsourcing.ToCloudEvent
	_ = cloudoutbox.ToCloudEvent
	_ = cloudqueue.ToCloudEvent
	_ = cloudrabbitstream.Encode
	_ = cloudtelemetry.InjectTraceContext
	_ = cloudworkflow.ToCloudEvent
}
