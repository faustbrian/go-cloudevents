# CloudEvents audit adapter

This module maps a deliberately small trusted subset of Golib audit metadata
through CloudEvents. It never treats an interoperability event as an audit
record and returns explicit losses for audit-owned fields.

Install the published `adapters/audit/v2.0.0` release with
`go get github.com/faustbrian/go-cloudevents/adapters/audit/v2@v2.0.0`.
Use `AddMetadata` for outbound metadata and `ExtractMetadata` only after the
application has made its trust decision. The caller retains audit-log,
authorization, transport, and lifecycle ownership.

## V2 migration

Move audit-adapter imports and `Metadata.Tenant` values to their `/v2`
namespaces together. The tenant field is the nominal
`github.com/faustbrian/go-tenancy/v2.TenantID`; trusted extraction uses the
public Tenancy adapter v2, and wrapped tenant errors match Tenancy v2 sentinels.
Audit core records, CloudEvents, explicit trust, input ownership and losses
retain their existing contracts. Published adapter v1 remains independently
available with its original tenant types and dependencies.

Source stays in `adapters/audit/` on main; no root-module major or
version-specific directory is introduced. `api/baseline.txt` preserves v1;
`api/v2-baseline.txt` describes the new major. Released facade v1 and v2 keep
their original contracts; facade-v3 source and target integration now consume
this public v2 adapter.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/audit/v2),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
