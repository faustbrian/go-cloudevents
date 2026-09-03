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

For shared package families, selection guidance, construction, ownership, and
lifecycle vocabulary, see the versioned [Golib ecosystem
index](https://github.com/faustbrian/go-library-tools/blob/v1.3.0/docs/ecosystem/README.md).

## Guarantees and limitations

The [complete guide](docs/reference.md) defines ownership, failure semantics,
bounds, concurrency, security, and unsupported behavior. Do not infer
additional guarantees beyond the documented module boundary.

## Documentation

- [Documentation index](docs/README.md)
- [Complete technical guide](docs/reference.md)
- [Go API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/golib)
- [Parent package documentation](../../docs/README.md)

## Compatibility and support

This module follows Semantic Versioning. Report vulnerabilities through the
[parent security policy](../../SECURITY.md).

## License

MIT. See [LICENSE](LICENSE).
