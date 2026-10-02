# CloudEvents schema-registry adapter

This module provides caller-invoked, bounded schema-registry resolution for
CloudEvents JSON Schema validation. Event decoding never performs lookup.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/schema-registry/v2@v2`.
Construct `JSONSchemaValidator` with an immutable URI allowlist, resolution
cache, adapter, availability policy, and positive timeout. The caller owns
registry credentials and lifecycle.

## Migration from v1

Version 2 uses `github.com/faustbrian/go-schema-registry/v2@v2.0.0`.
Update this adapter's import to
`github.com/faustbrian/go-cloudevents/adapters/schema-registry/v2` and use the
registry `/v2` imports, including `/v2/formats/jsonschema`, for the cache,
lookups, availability policy, and JSON Schema adapter supplied to
`JSONSchemaConfig`. Their v1 and v2 Go types are not interchangeable.

The CloudEvents root module stays on v1. The `SchemaValidator` interface,
CloudEvents error sentinels, exact trusted-URI mappings, caller-owned limits,
availability policy, and timeout behavior are unchanged. Registry v2 preserves
error classifications while keeping JSON payload diagnostics private.

The deprecated `adapters/golib` facade remains on its released v1 registry
adapter and types; it does not silently opt callers into v2. Migrate by
selecting this target adapter directly. Both module generations can coexist.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/schema-registry/v2),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
