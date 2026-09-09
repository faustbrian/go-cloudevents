package golib_test

import (
	"context"
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/faustbrian/go-cloudevents"
	golib "github.com/faustbrian/go-cloudevents/adapters/golib"
	"github.com/faustbrian/go-kafka"
	schemaregistry "github.com/faustbrian/go-schema-registry"
	registryjsonschema "github.com/faustbrian/go-schema-registry/formats/jsonschema"
)

const recipeSchemaURI = "https://schemas.example/orders/created/v1"

var (
	errRecipeKafkaClosed = errors.New("recipe kafka runtime: closed")
	errRecipeKafkaEmpty  = errors.New("recipe kafka runtime: no record")
)

// Example_kafkaSchemaCloudEventFlow is a non-releasable application
// composition. The application owns schema selection, Kafka settlement, and
// shutdown; the adapter only maps and validates immutable values.
func Example_kafkaSchemaCloudEventFlow() {
	operationCtx, cancelOperation := context.WithTimeout(context.Background(), time.Second)
	defer cancelOperation()
	shutdownCtx, cancelShutdown := context.WithTimeout(context.Background(), time.Second)
	defer cancelShutdown()

	runtime := &recipeKafkaRuntime{}
	handledID := ""
	err := runKafkaSchemaCloudEventRecipe(
		operationCtx,
		shutdownCtx,
		runtime,
		[]byte(`{"order_id":"A-123"}`),
		func(event cloudevents.Event) error {
			handledID = event.ID()
			return nil
		},
	)
	if err != nil {
		panic(err)
	}
	fmt.Println(
		handledID,
		runtime.committed.topic,
		runtime.committed.partition,
		runtime.committed.offset,
		runtime.shutdownOrder,
	)
	// Output: event-1 orders.created.v1 0 1 [consumer producer]
}

func TestKafkaSchemaCloudEventRecipeRejectsInvalidEventBeforePublish(t *testing.T) {
	t.Parallel()

	runtime := &recipeKafkaRuntime{}
	err := runKafkaSchemaCloudEventRecipe(
		context.Background(), context.Background(), runtime, []byte(`{}`),
		func(cloudevents.Event) error { return nil },
	)
	if !errors.Is(err, golib.ErrSchemaViolation) {
		t.Fatalf("recipe error = %v, want schema violation", err)
	}
	if len(runtime.records) != 0 {
		t.Fatalf("published records = %d, want 0", len(runtime.records))
	}
	if fmt.Sprint(runtime.shutdownOrder) != "[consumer producer]" {
		t.Fatalf("shutdown order = %v", runtime.shutdownOrder)
	}
}

func TestKafkaSchemaCloudEventRecipeDoesNotAcknowledgeHandlerFailure(t *testing.T) {
	t.Parallel()

	handlerErr := errors.New("projection unavailable")
	runtime := &recipeKafkaRuntime{}
	err := runKafkaSchemaCloudEventRecipe(
		context.Background(), context.Background(), runtime, []byte(`{"order_id":"A-123"}`),
		func(cloudevents.Event) error { return handlerErr },
	)
	if !errors.Is(err, handlerErr) {
		t.Fatalf("recipe error = %v, want handler cause", err)
	}
	if runtime.commitAttempted != nil || runtime.committed != nil {
		t.Fatal("handler failure reached offset commit")
	}
	if fmt.Sprint(runtime.shutdownOrder) != "[consumer producer]" {
		t.Fatalf("shutdown order = %v", runtime.shutdownOrder)
	}
}

func TestKafkaSchemaCloudEventRecipePreservesCommitFailureWithoutClaimingSettlement(t *testing.T) {
	t.Parallel()

	commitErr := errors.New("commit outcome unknown")
	runtime := &recipeKafkaRuntime{commitErr: commitErr}
	err := runKafkaSchemaCloudEventRecipe(
		context.Background(), context.Background(), runtime, []byte(`{"order_id":"A-123"}`),
		func(cloudevents.Event) error { return nil },
	)
	if !errors.Is(err, commitErr) {
		t.Fatalf("recipe error = %v, want commit cause", err)
	}
	want := recipeSettlement{topic: "orders.created.v1", partition: 0, offset: 1}
	if runtime.commitAttempted == nil || *runtime.commitAttempted != want {
		t.Fatalf("commit attempt = %#v, want %#v", runtime.commitAttempted, want)
	}
	if runtime.committed != nil {
		t.Fatalf("failed commit reported settlement = %#v", runtime.committed)
	}
}

func TestKafkaSchemaCloudEventRecipeJoinsOperationAndOrderedShutdownFailures(t *testing.T) {
	t.Parallel()

	publishErr := errors.New("publish outcome unknown")
	consumerErr := errors.New("consumer shutdown incomplete")
	producerErr := errors.New("producer drain incomplete")
	runtime := &recipeKafkaRuntime{
		publishErr:          publishErr,
		consumerShutdownErr: consumerErr,
		producerShutdownErr: producerErr,
	}
	err := runKafkaSchemaCloudEventRecipe(
		context.Background(), context.Background(), runtime, []byte(`{"order_id":"A-123"}`),
		func(cloudevents.Event) error { return nil },
	)
	if !errors.Is(err, publishErr) || !errors.Is(err, consumerErr) || !errors.Is(err, producerErr) {
		t.Fatalf("recipe error = %v, want operation and both shutdown causes", err)
	}
	if fmt.Sprint(runtime.shutdownOrder) != "[consumer producer]" {
		t.Fatalf("shutdown order = %v", runtime.shutdownOrder)
	}
}

func runKafkaSchemaCloudEventRecipe(
	operationCtx context.Context,
	shutdownCtx context.Context,
	runtime *recipeKafkaRuntime,
	payload []byte,
	handle func(cloudevents.Event) error,
) (resultErr error) {
	defer func() {
		resultErr = errors.Join(resultErr, runtime.Shutdown(shutdownCtx))
	}()

	validator, err := newRecipeRegistryValidator(operationCtx)
	if err != nil {
		return err
	}
	data, err := cloudevents.NewJSONData(payload)
	if err != nil {
		return err
	}
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID:              "event-1",
		Source:          "/orders",
		Type:            "orders.created",
		DataContentType: "application/json",
		DataSchema:      recipeSchemaURI,
	}, data)
	if err != nil {
		return err
	}
	if err := cloudevents.ValidateSchema(operationCtx, event, validator); err != nil {
		return err
	}
	record, err := golib.EncodeKafka(event, cloudevents.BinaryMode, golib.KafkaTransport{
		Topic: "orders.created.v1",
		Key:   []byte("A-123"),
	})
	if err != nil {
		return err
	}
	if err := runtime.Publish(operationCtx, record); err != nil {
		// A publish timeout or transport failure can have an unknown durable
		// outcome. Reconcile it instead of blindly retrying.
		return err
	}
	handler := kafka.HandlerFunc(func(ctx context.Context, consumed kafka.ConsumedRecord) error {
		message, _, err := golib.DecodeKafka(consumed, cloudevents.DefaultLimits())
		if err != nil {
			return err
		}
		if err := cloudevents.ValidateSchema(ctx, message.Event, validator); err != nil {
			return err
		}
		// Returning an error leaves the borrowed record unsettled. A real
		// kafka.Consumer owns offset commit and retry/dead-letter policy.
		return handle(message.Event)
	})
	return runtime.RunOnce(operationCtx, handler)
}

func newRecipeRegistryValidator(ctx context.Context) (golib.RegistryJSONSchemaValidator, error) {
	adapter, err := registryjsonschema.New(registryjsonschema.Config{
		MaxSchemaBytes: 1024, MaxTotalSchemaBytes: 1024,
		MaxPayloadBytes: 1024, MaxResources: 1,
	})
	if err != nil {
		return golib.RegistryJSONSchemaValidator{}, err
	}
	schema, err := schemaregistry.Compile(ctx, schemaregistry.Definition{
		Format:  schemaregistry.FormatJSONSchema,
		Content: []byte(`{"type":"object","required":["order_id"],"properties":{"order_id":{"type":"string"}},"additionalProperties":false}`),
	}, adapter)
	if err != nil {
		return golib.RegistryJSONSchemaValidator{}, err
	}
	subject := schemaregistry.Subject{Name: "orders.created.v1"}
	lookup := schemaregistry.Latest(subject)
	cache, err := schemaregistry.NewResolveCache(
		recipeSchemaResolver{result: schemaregistry.ResolveResult{
			Schema: schema,
			ID: schemaregistry.ProviderID{
				Provider: "recipe",
				Value:    "orders-created-v1",
			},
			Subject: subject, Version: schemaregistry.Version{Number: 1},
			Lifecycle: schemaregistry.LifecycleAvailable,
		}},
		schemaregistry.ResolveCacheConfig{
			MaxEntries: 1, MaxConcurrent: 1,
			FreshFor: time.Minute, StaleFor: time.Minute, NegativeFor: time.Minute,
			Clock: recipeClock{now: time.Unix(0, 0)},
		},
	)
	if err != nil {
		return golib.RegistryJSONSchemaValidator{}, err
	}
	return golib.NewRegistryJSONSchemaValidator(golib.RegistryJSONSchemaConfig{
		Cache: cache, SchemaLookups: map[string]schemaregistry.Lookup{recipeSchemaURI: lookup},
		Adapter: adapter, AvailabilityPolicy: schemaregistry.FailClosed, Timeout: time.Second,
	})
}

type recipeSchemaResolver struct {
	result schemaregistry.ResolveResult
}

func (resolver recipeSchemaResolver) Resolve(
	ctx context.Context,
	_ schemaregistry.Lookup,
) (schemaregistry.ResolveResult, error) {
	if err := ctx.Err(); err != nil {
		return schemaregistry.ResolveResult{}, err
	}
	return resolver.result, nil
}

type recipeClock struct {
	now time.Time
}

func (clock recipeClock) Now() time.Time { return clock.now }

type recipeKafkaRuntime struct {
	records             []kafka.ProducerRecord
	publishErr          error
	commitErr           error
	consumerShutdownErr error
	producerShutdownErr error
	commitAttempted     *recipeSettlement
	committed           *recipeSettlement
	shutdownOrder       []string
}

func (runtime *recipeKafkaRuntime) Publish(ctx context.Context, record kafka.ProducerRecord) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(runtime.shutdownOrder) != 0 {
		return errRecipeKafkaClosed
	}
	if runtime.publishErr != nil {
		return runtime.publishErr
	}
	runtime.records = append(runtime.records, cloneRecipeRecord(record))
	return nil
}

func (runtime *recipeKafkaRuntime) RunOnce(ctx context.Context, handler kafka.Handler) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(runtime.records) == 0 {
		return errRecipeKafkaEmpty
	}
	record := runtime.records[0]
	consumed := kafka.ConsumedRecord{
		Topic: record.Topic, Key: record.Key, Value: record.Value, Headers: record.Headers,
		Timestamp: record.Timestamp, TimestampType: kafka.TimestampCreateTime,
		Partition: 0, Offset: 1,
	}
	if err := handler.Handle(ctx, consumed); err != nil {
		return err
	}
	settlement := recipeSettlement{
		topic: consumed.Topic, partition: consumed.Partition, offset: consumed.Offset,
	}
	runtime.commitAttempted = &settlement
	if runtime.commitErr != nil {
		return runtime.commitErr
	}
	runtime.committed = &settlement
	return nil
}

func (runtime *recipeKafkaRuntime) Shutdown(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if len(runtime.shutdownOrder) == 0 {
		runtime.shutdownOrder = append(runtime.shutdownOrder, "consumer")
		consumerErr := runtime.consumerShutdownErr
		runtime.shutdownOrder = append(runtime.shutdownOrder, "producer")
		return errors.Join(consumerErr, runtime.producerShutdownErr)
	}
	return nil
}

type recipeSettlement struct {
	topic     string
	partition int32
	offset    int64
}

func cloneRecipeRecord(record kafka.ProducerRecord) kafka.ProducerRecord {
	clone := record
	clone.Key = append([]byte(nil), record.Key...)
	clone.Value = append([]byte(nil), record.Value...)
	clone.Headers = make([]kafka.Header, len(record.Headers))
	for index, header := range record.Headers {
		clone.Headers[index] = kafka.Header{
			Key: header.Key, Value: append([]byte(nil), header.Value...),
		}
	}
	return clone
}
