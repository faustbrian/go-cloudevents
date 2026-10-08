# Golib CloudEvents adapters

`golib` is the optional integration module between the transport-independent
CloudEvents package and Golib's canonical event, transport, workflow,
metadata, audit, and schema contracts. Importing it performs no registration,
network access, schema lookup, telemetry emission, or background work.

This is a deprecated compatibility facade. Version 4 adopts public JSON Schema
v2, Outbox v2 and Schema Registry v3 alongside Tenancy v2, EventSourcing v2 and
Workflow v2; released v1, v2 and v3 retain their original types. All generations are
excluded from the recommended set because the broad bridge
owns 26 module dependencies. New code and migrations should import the target-oriented
module listed in the parent [adoption guide](https://github.com/faustbrian/go-cloudevents#adoption-guidance)
so they carry only the relevant optional dependencies.

Conversions retain canonical state that CloudEvents cannot represent and
return explicit loss reports. Queue and outbox conversions are Golib mappings,
not official CloudEvents protocol bindings. Schema resolution occurs only when
the caller explicitly invokes CloudEvents schema validation with a configured
registry validator.

## Install

```sh
go get github.com/faustbrian/go-cloudevents/adapters/golib/v4@v4.0.0
```

Requires Go 1.27. This command applies after facade `adapters/golib/v4.0.0`
publication.

## Migration from v3

Update facade imports to
`github.com/faustbrian/go-cloudevents/adapters/golib/v4`. Move JSON Schema imports
to `github.com/faustbrian/go-json-schema/v2`, outbox imports to
`github.com/faustbrian/go-transactional-outbox/v2`, and registry imports (including
format adapters) to `github.com/faustbrian/go-schema-registry/v3`.
Direct validators use JSON Schema2 compiled schemas, outbox conversions accept
and return Outbox2 envelopes, and registry validator configuration uses Registry3
cache, lookup and format-adapter types. These nominal types are not
interchangeable with their legacy identities; update paired imports together.

Tenancy, EventSourcing and Workflow retain their already adopted v2 identities.
Applications upgrading from facade v1 or v2 must also use public Tenancy2,
EventSourcing2 and Workflow2 producer imports.

Stable IDs, sequence, occurrence time, copied payloads, nil-versus-empty data,
retained state, explicit trust decisions and losses keep their existing
semantics. The facade consumes the actual published JSONSchema2, Outbox2 and
Schema Registry3 target adapters alongside Audit2, Queue2, EventSourcing2,
Tenancy2 and Workflow2. The CloudEvents root
remains a v1 module, selected at v1.1.1; this migration does not create a root
major or change unrelated transport or telemetry contracts.

Source remains in `adapters/golib/` on main, with independent tag
`adapters/golib/v4.0.0`. No root-module release or version-specific source
directory is required. `api/v1.1.1.txt` preserves the exact released facade API
from commit `5559c521abeb0ac979a06a9405185b8486033c88`;
`api/baseline.txt` preserves the released facade-v2 API and
`api/v3-baseline.txt` preserves the released v3 API;
`api/v4-baseline.txt` describes the new current major. Prefer importing only the
target-oriented adapters your application needs.

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

[`Example_kafkaSchemaCloudEventFlow`](../../integration/target-adapters/kafka_schema_cloudevents_example_test.go)
is the executable, non-releasable reference composition for CloudEvents,
Kafka, JSON Schema, and schema registry. It lives in the target-adapter
integration module, imports the three target adapters directly, and keeps
orchestration in the application:

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

The target adapters own only mapping and validation. The application owns
registry credentials and availability policy, Kafka brokers, topics, retry and
dead-letter policy, acknowledgements, correlation and telemetry, and the
runtime lifecycle. Optional provider and telemetry adapters remain application
choices.

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
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/golib/v4)
- [Parent package documentation](https://github.com/faustbrian/go-cloudevents/tree/main/docs)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](https://github.com/faustbrian/go-cloudevents/security/policy).

## License

MIT. See [LICENSE](LICENSE).
