# CloudEvents queue adapter

This module maps Golib queue jobs to CloudEvents while retaining retry,
settlement, execution, and operational state outside the event.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/queue@v1`.
Use `ToCloudEvent` and retain the returned job; pass it to `FromCloudEvent` to
reconstruct the canonical queue value. The application owns queue I/O,
acknowledgements, retry policy, and lifecycle.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/queue),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
