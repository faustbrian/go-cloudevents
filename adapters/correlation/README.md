# CloudEvents correlation adapter

This module maps caller-owned Golib correlation identifiers through
CloudEvents extensions without hidden I/O or mutable aliases.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/correlation@v1`.
Use `Add` for outbound identifiers and `Extract` only after the application
has made its trust decision. The caller owns identifier policy and lifecycle.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/correlation),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
