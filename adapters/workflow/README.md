# CloudEvents workflow adapter

This module maps Golib durable workflow history to CloudEvents while retaining
workflow-owned sequence, definition, retry, compensation, and scheduling
state outside the event.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/workflow@v1`.
Use `ToCloudEvent` and retain its `State`; pass that state to `FromCloudEvent`
when reconstructing the canonical history event. Runtime lifecycle is external.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/workflow),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
