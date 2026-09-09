# CloudEvents tenancy adapter

This module maps validated Golib tenant routing identity through one
CloudEvents extension. It does not authenticate or authorize the tenant.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/tenancy@v1`.
Use `Add` for outbound identity and `Extract` only after the application has
made its trust decision. Authorization and transport routing remain external.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/tenancy),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
