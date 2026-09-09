# CloudEvents outbox adapter

This module maps Golib transactional-outbox envelopes to CloudEvents while
retaining relay-owned state outside the interoperability event.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/outbox@v1`.
Use `ToCloudEvent` and retain the returned envelope; pass it to
`FromCloudEvent` to reconstruct the canonical outbox value. The application
owns transactions, persistence, dispatch, retries, and settlement.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/outbox),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
