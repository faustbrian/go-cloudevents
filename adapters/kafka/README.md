# CloudEvents Kafka adapter

This module maps CloudEvents through Golib Kafka records. It performs no
broker I/O and owns no producer, consumer, offset, retry, or shutdown state.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/kafka@v1`.
Use `Encode` with an explicit CloudEvents content mode and transport metadata;
use `Decode` with caller-selected limits. Returned bytes do not alias input.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/kafka),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
