package schemaregistry_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/faustbrian/go-cloudevents"
	cloudregistry "github.com/faustbrian/go-cloudevents/adapters/schema-registry"
	golibregistry "github.com/faustbrian/go-schema-registry"
	registryjsonschema "github.com/faustbrian/go-schema-registry/formats/jsonschema"
)

type resolverFunc func(context.Context, golibregistry.Lookup) (golibregistry.ResolveResult, error)

func (fn resolverFunc) Resolve(ctx context.Context, lookup golibregistry.Lookup) (golibregistry.ResolveResult, error) {
	return fn(ctx, lookup)
}

type fixedClock struct{ now time.Time }

func (clock fixedClock) Now() time.Time { return clock.now }

type passthroughCanonicalizer struct{}

func (passthroughCanonicalizer) Canonicalize(_ context.Context, definition golibregistry.Definition) ([]byte, error) {
	return append([]byte(nil), definition.Content...), nil
}

func TestNewJSONSchemaValidatorRejectsUnsafeConfiguration(t *testing.T) {
	cache := newCache(t, resolverFunc(func(context.Context, golibregistry.Lookup) (golibregistry.ResolveResult, error) {
		return golibregistry.ResolveResult{}, errors.New("must not resolve")
	}))
	adapter := newAdapter(t)
	lookup := schemaLookup()
	valid := cloudregistry.JSONSchemaConfig{
		Cache: cache, SchemaLookups: map[string]golibregistry.Lookup{"schema": lookup}, Adapter: adapter,
		AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second,
	}

	tests := map[string]cloudregistry.JSONSchemaConfig{
		"nil cache":          {SchemaLookups: valid.SchemaLookups, Adapter: adapter, AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second},
		"empty mappings":     {Cache: cache, Adapter: adapter, AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second},
		"empty URI":          {Cache: cache, SchemaLookups: map[string]golibregistry.Lookup{"": lookup}, Adapter: adapter, AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second},
		"empty lookup":       {Cache: cache, SchemaLookups: map[string]golibregistry.Lookup{"schema": {}}, Adapter: adapter, AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second},
		"nil adapter":        {Cache: cache, SchemaLookups: valid.SchemaLookups, AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second},
		"implicit policy":    {Cache: cache, SchemaLookups: valid.SchemaLookups, Adapter: adapter, Timeout: time.Second},
		"unsupported policy": {Cache: cache, SchemaLookups: valid.SchemaLookups, Adapter: adapter, AvailabilityPolicy: golibregistry.AvailabilityPolicy("unsafe"), Timeout: time.Second},
		"zero timeout":       {Cache: cache, SchemaLookups: valid.SchemaLookups, Adapter: adapter, AvailabilityPolicy: golibregistry.FailClosed},
		"negative timeout":   {Cache: cache, SchemaLookups: valid.SchemaLookups, Adapter: adapter, AvailabilityPolicy: golibregistry.FailClosed, Timeout: -time.Second},
	}
	for name, config := range tests {
		t.Run(name, func(t *testing.T) {
			if _, err := cloudregistry.NewJSONSchemaValidator(config); !errors.Is(err, cloudregistry.ErrSchemaMapping) {
				t.Fatalf("configuration error = %v, want %v", err, cloudregistry.ErrSchemaMapping)
			}
		})
	}

	for _, policy := range []golibregistry.AvailabilityPolicy{
		golibregistry.FailClosed,
		golibregistry.AllowStale,
		golibregistry.CacheOnly,
		golibregistry.ReturnUnavailable,
	} {
		config := valid
		config.AvailabilityPolicy = policy
		if _, err := cloudregistry.NewJSONSchemaValidator(config); err != nil {
			t.Fatalf("supported policy %q: %v", policy, err)
		}
	}
}

func TestJSONSchemaValidatorRejectsUntrustedInputsBeforeResolution(t *testing.T) {
	called := false
	validator := newValidator(t, resolverFunc(func(context.Context, golibregistry.Lookup) (golibregistry.ResolveResult, error) {
		called = true
		return golibregistry.ResolveResult{}, errors.New("must not resolve")
	}))

	//lint:ignore SA1012 Nil is the contract under test.
	if err := validator.Validate(nil, "schema", "application/json", []byte(`{}`)); !errors.Is(err, cloudevents.ErrContextRequired) { //nolint:staticcheck // Nil is the contract under test.
		t.Fatalf("nil context error = %v, want %v", err, cloudevents.ErrContextRequired)
	}
	if err := (cloudregistry.JSONSchemaValidator{}).Validate(context.Background(), "schema", "application/json", []byte(`{}`)); !errors.Is(err, cloudregistry.ErrSchemaMapping) {
		t.Fatalf("zero-value error = %v, want %v", err, cloudregistry.ErrSchemaMapping)
	}
	if err := validator.Validate(context.Background(), "schema", "text/plain", []byte(`{}`)); !errors.Is(err, cloudregistry.ErrSchemaMapping) {
		t.Fatalf("content type error = %v, want %v", err, cloudregistry.ErrSchemaMapping)
	}
	if err := validator.Validate(context.Background(), "untrusted", "application/json", []byte(`{}`)); !errors.Is(err, cloudregistry.ErrSchemaMapping) {
		t.Fatalf("unmapped URI error = %v, want %v", err, cloudregistry.ErrSchemaMapping)
	}
	if called {
		t.Fatal("resolver called for rejected input")
	}
}

func TestJSONSchemaValidatorPreservesResolutionAndFormatFailures(t *testing.T) {
	wantResolutionErr := errors.New("registry unavailable")
	validator := newValidator(t, resolverFunc(func(context.Context, golibregistry.Lookup) (golibregistry.ResolveResult, error) {
		return golibregistry.ResolveResult{}, wantResolutionErr
	}))
	if err := validator.Validate(context.Background(), "schema", "application/json", []byte(`{}`)); !errors.Is(err, wantResolutionErr) {
		t.Fatalf("resolution error = %v, want %v", err, wantResolutionErr)
	}

	avro, err := golibregistry.Compile(context.Background(), golibregistry.Definition{
		Format: golibregistry.FormatAvro, Content: []byte(`{"type":"record","name":"Event","fields":[]}`),
	}, passthroughCanonicalizer{})
	if err != nil {
		t.Fatal(err)
	}
	validator = newValidator(t, resolverFunc(func(context.Context, golibregistry.Lookup) (golibregistry.ResolveResult, error) {
		return schemaResult(avro), nil
	}))
	if err := validator.Validate(context.Background(), "schema", "application/json", []byte(`{}`)); !errors.Is(err, cloudregistry.ErrSchemaMapping) {
		t.Fatalf("registry format error = %v, want %v", err, cloudregistry.ErrSchemaMapping)
	}
}

func TestJSONSchemaValidatorReportsPayloadOutcomes(t *testing.T) {
	adapter := newAdapter(t)
	schema := compileJSONSchema(t, adapter)
	lookup := schemaLookup()
	cache := newCache(t, resolverFunc(func(_ context.Context, got golibregistry.Lookup) (golibregistry.ResolveResult, error) {
		if got != lookup {
			return golibregistry.ResolveResult{}, errors.New("unexpected lookup")
		}
		return schemaResult(schema), nil
	}))
	validator, err := cloudregistry.NewJSONSchemaValidator(cloudregistry.JSONSchemaConfig{
		Cache: cache, SchemaLookups: map[string]golibregistry.Lookup{"schema": lookup}, Adapter: adapter,
		AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	if err := validator.Validate(context.Background(), "schema", "application/problem+json", []byte(`{"id":1}`)); err != nil {
		t.Fatalf("valid payload: %v", err)
	}
	if err := validator.Validate(context.Background(), "schema", "application/json", []byte(`{}`)); !errors.Is(err, cloudregistry.ErrSchemaViolation) {
		t.Fatalf("schema violation error = %v, want %v", err, cloudregistry.ErrSchemaViolation)
	}

	cancelCtx, cancel := context.WithCancel(context.Background())
	cancelingCache := newCache(t, resolverFunc(func(context.Context, golibregistry.Lookup) (golibregistry.ResolveResult, error) {
		cancel()
		return schemaResult(schema), nil
	}))
	canceling, err := cloudregistry.NewJSONSchemaValidator(cloudregistry.JSONSchemaConfig{
		Cache: cancelingCache, SchemaLookups: map[string]golibregistry.Lookup{"schema": lookup}, Adapter: adapter,
		AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	if err := canceling.Validate(cancelCtx, "schema", "application/json", []byte(`{"id":1}`)); !errors.Is(err, context.Canceled) {
		t.Fatalf("canceled validation error = %v, want %v", err, context.Canceled)
	}
}

func TestJSONSchemaValidatorOwnsLookupSnapshot(t *testing.T) {
	adapter := newAdapter(t)
	schema := compileJSONSchema(t, adapter)
	lookup := schemaLookup()
	source := map[string]golibregistry.Lookup{"schema": lookup}
	cache := newCache(t, resolverFunc(func(_ context.Context, got golibregistry.Lookup) (golibregistry.ResolveResult, error) {
		if got != lookup {
			return golibregistry.ResolveResult{}, errors.New("unexpected lookup")
		}
		return schemaResult(schema), nil
	}))
	validator, err := cloudregistry.NewJSONSchemaValidator(cloudregistry.JSONSchemaConfig{
		Cache: cache, SchemaLookups: source, Adapter: adapter,
		AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}

	source["schema"] = golibregistry.ByProviderID(golibregistry.ProviderID{Provider: "attacker", Value: "replacement"})
	source["untrusted"] = lookup
	if err := validator.Validate(context.Background(), "schema", "application/json", []byte(`{"id":1}`)); err != nil {
		t.Fatalf("validation from owned lookup snapshot: %v", err)
	}
	if err := validator.Validate(context.Background(), "untrusted", "application/json", []byte(`{"id":1}`)); !errors.Is(err, cloudregistry.ErrSchemaMapping) {
		t.Fatalf("post-construction map injection error = %v, want %v", err, cloudregistry.ErrSchemaMapping)
	}
}

func newValidator(t *testing.T, resolver golibregistry.Resolver) cloudregistry.JSONSchemaValidator {
	t.Helper()
	validator, err := cloudregistry.NewJSONSchemaValidator(cloudregistry.JSONSchemaConfig{
		Cache: newCache(t, resolver), SchemaLookups: map[string]golibregistry.Lookup{"schema": schemaLookup()}, Adapter: newAdapter(t),
		AvailabilityPolicy: golibregistry.FailClosed, Timeout: time.Second,
	})
	if err != nil {
		t.Fatal(err)
	}
	return validator
}

func newCache(t *testing.T, resolver golibregistry.Resolver) *golibregistry.ResolveCache {
	t.Helper()
	cache, err := golibregistry.NewResolveCache(resolver, golibregistry.ResolveCacheConfig{
		MaxEntries: 4, MaxConcurrent: 1, FreshFor: time.Minute, StaleFor: time.Minute,
		NegativeFor: time.Minute, Clock: fixedClock{now: time.Unix(1, 0)},
	})
	if err != nil {
		t.Fatal(err)
	}
	return cache
}

func newAdapter(t *testing.T) *registryjsonschema.Adapter {
	t.Helper()
	adapter, err := registryjsonschema.New(registryjsonschema.Config{
		MaxSchemaBytes: 1024, MaxTotalSchemaBytes: 2048, MaxPayloadBytes: 1024, MaxResources: 1,
	})
	if err != nil {
		t.Fatal(err)
	}
	return adapter
}

func compileJSONSchema(t *testing.T, adapter *registryjsonschema.Adapter) golibregistry.Schema {
	t.Helper()
	schema, err := golibregistry.Compile(context.Background(), golibregistry.Definition{
		Format: golibregistry.FormatJSONSchema, Content: []byte(`{"type":"object","required":["id"]}`),
	}, adapter)
	if err != nil {
		t.Fatal(err)
	}
	return schema
}

func schemaLookup() golibregistry.Lookup {
	return golibregistry.ByProviderID(golibregistry.ProviderID{Provider: "test", Scope: "local", Value: "orders-v1"})
}

func schemaResult(schema golibregistry.Schema) golibregistry.ResolveResult {
	return golibregistry.ResolveResult{
		Schema: schema, ID: schemaLookup().ProviderID(), Lifecycle: golibregistry.LifecycleAvailable,
	}
}
