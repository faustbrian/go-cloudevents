# CloudEvents schema-registry adapter

[Documentation index](../../docs/README.md) · [Repository entry](../../README.md)

This module provides caller-invoked, bounded schema-registry resolution for
CloudEvents JSON Schema validation. Event decoding never performs lookup.

This independently versioned adapter is published at `v3.0.0`. Install
with `go get github.com/faustbrian/go-cloudevents/adapters/schema-registry/v3@v3.0.0`.
Construct `JSONSchemaValidator` with an immutable URI allowlist, resolution
cache, adapter, availability policy, and positive timeout. The caller owns
registry credentials and lifecycle.

## Migration from v2

Version 3 uses `github.com/faustbrian/go-schema-registry/v3@v3.0.0`.
Update this adapter's import to
`github.com/faustbrian/go-cloudevents/adapters/schema-registry/v3` and use the
registry `/v3` imports, including `/v3/formats/jsonschema`, for the cache,
lookups, availability policy, and JSON Schema adapter supplied to
`JSONSchemaConfig`. Their v2 and v3 Go types are not interchangeable.

The CloudEvents root module stays on v1. The `SchemaValidator` interface,
CloudEvents error sentinels, exact trusted-URI mappings, caller-owned limits,
availability policy, and timeout behavior are unchanged. Registry v3 preserves
error classifications while keeping JSON payload diagnostics private.

The deprecated Golib facade v3 retains its released registry-adapter v2 and
types. Published facade v4.0.0 adopts the JSON Schema, Registry and Outbox
target majors together. Prefer this target adapter directly when only this
boundary is needed; existing public module generations can coexist.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/schema-registry/v3),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).

Choose the next guide in the [documentation index](../../docs/README.md)
or return to the [repository entry](../../README.md) for installation and ownership.
