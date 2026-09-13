# CloudEvents audit adapter

This module maps a deliberately small trusted subset of Golib audit metadata
through CloudEvents. It never treats an interoperability event as an audit
record and returns explicit losses for audit-owned fields.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/audit@v1`.
Use `AddMetadata` for outbound metadata and `ExtractMetadata` only after the
application has made its trust decision. The caller retains audit-log,
authorization, transport, and lifecycle ownership.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/audit),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
