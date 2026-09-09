# Golib CloudEvents adapters

`golib` is the optional integration module between the transport-independent
CloudEvents package and Golib's canonical event, transport, workflow,
metadata, audit, and schema contracts. Importing it performs no registration,
network access, schema lookup, telemetry emission, or background work.

Conversions retain canonical state that CloudEvents cannot represent and
return explicit loss reports. Queue and outbox conversions are Golib mappings,
not official CloudEvents protocol bindings. Schema resolution occurs only when
the caller explicitly invokes CloudEvents schema validation with a configured
registry validator.

## Install

```sh
go get github.com/faustbrian/go-cloudevents/adapters/golib@v1
```

## Quick start

```go
message := job.Message{
    Timeout: time.Minute,
    Body: []byte(`{"order":"A-123"}`),
    Metadata: &job.Metadata{
        OriginalID: "job-1",
        JobType: "order.notify",
        ContentType: "application/json",
    },
}

event, retained, report, err := golib.QueueToCloudEvent(
    message,
    golib.QueueOptions{Source: "/queue/orders"},
)
```

The compiling examples in this module contain complete imports and setup.

### Kafka event flow with a registry schema

[`Example_kafkaSchemaCloudEventFlow`](kafka_schema_cloudevents_example_test.go)
is the executable, non-releasable reference composition for CloudEvents,
Kafka, JSON Schema, and schema registry. It uses only public package APIs and
keeps orchestration in the application:

1. The application selects the schema URI and bounded registry lookup, then
   constructs the JSON Schema adapter, resolution cache, and validator.
2. The producer validates the immutable event before the adapter encodes a
   Kafka record. The application-owned Kafka producer takes ownership of the
   encoded bytes when it accepts the record.
3. A public `kafka.HandlerFunc` borrows the Kafka record, decodes it into an
   owned CloudEvent, resolves and validates its schema, and returns the
   application result. `kafka.Consumer` owns offset commits and commits only a
   successfully handled contiguous prefix.
4. Validation and handler failures therefore remain uncommitted. Publish and
   commit timeouts may have unknown outcomes and must be reconciled rather than
   blindly retried.
5. Shutdown uses its own bounded context: intake stops first, the consumer
   drains and closes, then the producer drains and closes even if consumer
   shutdown fails. Operation and cleanup failures are joined so no cause is
   discarded.

`go-cloudevents/adapters/golib` owns only mapping and validation. The
application owns registry credentials and availability policy, Kafka brokers,
topics, retry and dead-letter policy, acknowledgements, correlation and
telemetry, and the runtime lifecycle. The required module versions are the
independent versions declared in this adapter module's `go.mod`; optional
provider and telemetry adapters remain application choices.

For shared package families, selection guidance, construction, ownership, and
lifecycle vocabulary, see the versioned [Golib ecosystem
index](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/README.md)
and its [protocols-and-descriptions package
guidance](https://github.com/faustbrian/go-library-tools/blob/v1.4.0/docs/ecosystem/design-language.md#package-families-and-selection).

## Guarantees and limitations

The [complete guide](docs/reference.md) defines ownership, failure semantics,
bounds, concurrency, security, and unsupported behavior. Do not infer
additional guarantees beyond the documented module boundary.

## Documentation

- [Documentation index](docs/README.md)
- [Complete technical guide](docs/reference.md)
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/golib)
- [Parent package documentation](https://github.com/faustbrian/go-cloudevents/tree/main/docs)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](https://github.com/faustbrian/go-cloudevents/security/policy).

## License

MIT. See [LICENSE](LICENSE).
