# CloudEvents RabbitMQ Streams adapter

This module maps structured CloudEvents through Golib RabbitMQ Streams
messages without claiming a non-existent binary AMQP binding.

Install with `go get github.com/faustbrian/go-cloudevents/adapters/rabbitstream@v1`.
Use `Encode` with one explicit stream target and `Decode` with caller-selected
CloudEvents limits. Broker offsets, routing, resources, and lifecycle remain
caller-owned.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/rabbitstream),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
