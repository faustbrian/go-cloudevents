# CloudEvents JSON Schema adapter

This module adapts one caller-compiled Golib JSON Schema to the CloudEvents
`SchemaValidator` boundary. It performs no registry lookup or network I/O.

The v2 adapter is prepared in this source and pending publication. After release,
install with `go get github.com/faustbrian/go-cloudevents/adapters/jsonschema/v2@v2`.
Construct a `Validator` with the exact schema URI and compiled schema, then
pass it explicitly to CloudEvents validation. The caller owns compilation,
availability, and lifecycle.

Version 2 consumes `github.com/faustbrian/go-json-schema/v2` v2.0.0 compiled
schemas. Upgrade adapter and schema imports together; their v1 and v2 schema
types are not interchangeable. The CloudEvents root and its `SchemaValidator`
interface remain v1. Compilation policy, resource limits and remote-reference
admission remain explicit caller responsibilities. Parser, limit and context
errors are preserved, and schema violations retain the CloudEvents sentinel.

The released Golib facade v3 still uses the old target cohort. Its combined
v4 migration follows publication of the JSON Schema, Registry and Outbox target
adapters; selecting this adapter alone does not migrate the facade.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/jsonschema/v2),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
