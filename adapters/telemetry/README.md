# CloudEvents telemetry adapter

This module maps caller-owned Golib propagation policy through CloudEvents
trace extensions without owning telemetry initialization or shutdown.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/telemetry@v1`.
Use `InjectTraceContext` and `ExtractTraceContext` with an explicit policy and
trust decision. Baggage is not flattened and is reported as conversion loss.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/telemetry),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
