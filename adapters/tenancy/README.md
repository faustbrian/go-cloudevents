# CloudEvents tenancy adapter

This module maps validated Golib tenant routing identity through one
CloudEvents extension. It does not authenticate or authorize the tenant.

The v2 source uses `github.com/faustbrian/go-tenancy/v2` tenant identities.
After `adapters/tenancy/v2.0.0` is published, install with
`go get github.com/faustbrian/go-cloudevents/adapters/tenancy/v2@v2`.
Update both adapter and tenant imports together: v1 and v2 tenant types and
error sentinels have distinct identities. Released v1 consumers retain the
existing contract; the `tenantid` extension format is unchanged.
Use `Add` for outbound identity and `Extract` only after the application has
made its trust decision. Authorization and transport routing remain external.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/tenancy/v2),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
