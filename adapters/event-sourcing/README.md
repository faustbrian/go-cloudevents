# CloudEvents event-sourcing adapter

This module maps Golib event-sourcing messages to CloudEvents while keeping
stream versions, positions, persistence timestamps, and metadata in explicit
caller-owned state.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/event-sourcing@v1`.
Use `ToCloudEvent` and retain its `State`; pass that state to `FromCloudEvent`
when reconstructing the canonical message. No storage or transport is owned.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/event-sourcing),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
