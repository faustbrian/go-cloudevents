# CloudEvents JSON Schema adapter

This module adapts one caller-compiled Golib JSON Schema to the CloudEvents
`SchemaValidator` boundary. It performs no registry lookup or network I/O.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/jsonschema@v1`.
Construct a `Validator` with the exact schema URI and compiled schema, then
pass it explicitly to CloudEvents validation. The caller owns compilation,
availability, and lifecycle.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/jsonschema),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
