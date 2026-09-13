# CloudEvents schema-registry adapter

This module provides caller-invoked, bounded schema-registry resolution for
CloudEvents JSON Schema validation. Event decoding never performs lookup.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/schema-registry@v1`.
Construct `JSONSchemaValidator` with an immutable URI allowlist, resolution
cache, adapter, availability policy, and positive timeout. The caller owns
registry credentials and lifecycle.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/schema-registry),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
