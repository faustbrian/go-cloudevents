# CloudEvents outbox adapter

This module maps Golib transactional-outbox envelopes to CloudEvents while
retaining relay-owned state outside the interoperability event.

The v2 adapter is prepared in this source and pending publication. After release,
install with `go get github.com/faustbrian/go-cloudevents/adapters/outbox/v2@v2`.
Use `ToCloudEvent` and retain the returned envelope; pass it to
`FromCloudEvent` to reconstruct the canonical outbox value. The application
owns transactions, persistence, dispatch, retries, and settlement.

Version 2 consumes `github.com/faustbrian/go-transactional-outbox/v2` v2.0.0
envelopes. Upgrade adapter and outbox imports together; v1 and v2 envelopes are
distinct Go types. The CloudEvents root stays on v1. Conversion still clones
mutable state and preserves nil versus empty payloads, explicit losses and
identifier collision checks; it does not perform persistence or settlement.

The released Golib facade v3 still uses the old target cohort. Its combined v4
migration follows publication of all three schema/registry/outbox targets.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/outbox/v2),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
