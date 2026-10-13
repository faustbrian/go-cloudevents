# Package map

[Documentation index](README.md) · [Repository entry](../README.md)

This map follows the current source modules, not a synchronized release.
Use the repository README for published installation versions; an unpublished
source major is not installable merely because its package appears here.
Internal packages and non-releasable fixtures are intentionally excluded.

## github.com/faustbrian/go-cloudevents

`cloudevents` is Golib's transport-independent CloudEvents interoperability envelope. It implements the stable 1.0 information model, deterministic JSON event and batch serialization, and selected HTTP and Kafka mappings. It does not replace richer domain, event-sourcing, outbox, queue, workflow, audit, or application envelopes.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents`](../doc.go) | Use for CloudEvents information-model and wire interoperability; retain a richer domain envelope when it owns application, persistence, workflow, or settlement state. |

## github.com/faustbrian/go-cloudevents/adapters/golib/v4

`golib` is the deprecated CloudEvents compatibility facade adopting public Tenancy2, EventSourcing2, JSONSchema2, Outbox2 and Schema Registry3 alongside Workflow2. It retains explicit conversion, ownership, trust and loss contracts; importing it performs no registration, network access, schema lookup, telemetry emission or background work. Prefer independently versioned target adapters instead of its broad 26-dependency bridge.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/golib/v4`](../adapters/golib/doc.go) | Deprecated compatibility facade for explicit current public producer adoption; prefer the relevant target-oriented adapter from the parent adoption guide. |

## github.com/faustbrian/go-cloudevents/adapters/audit/v2

Map CloudEvents explicitly to and from Golib audit metadata while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/audit/v2`](../adapters/audit/audit.go) | Use when an application owns both a CloudEvent and Golib audit metadata contract. |

## github.com/faustbrian/go-cloudevents/adapters/correlation

Map CloudEvents explicitly to and from Golib correlation identifiers while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/correlation`](../adapters/correlation/correlation.go) | Use when an application owns both a CloudEvent and Golib correlation identifiers contract. |

## github.com/faustbrian/go-cloudevents/adapters/event-sourcing/v2

Map CloudEvents explicitly to and from Golib event-sourcing messages while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/event-sourcing/v2`](../adapters/event-sourcing/event_sourcing.go) | Use when an application owns both a CloudEvent and Golib event-sourcing v2 messages contract. |

## github.com/faustbrian/go-cloudevents/adapters/jsonschema/v2

Map CloudEvents explicitly to and from Golib compiled JSON Schema validation while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/jsonschema/v2`](../adapters/jsonschema/jsonschema.go) | Use when an application owns both a CloudEvent and Golib compiled JSON Schema validation contract. |

## github.com/faustbrian/go-cloudevents/adapters/kafka

Map CloudEvents explicitly to and from Golib Kafka records while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/kafka`](../adapters/kafka/kafka.go) | Use when an application owns both a CloudEvent and Golib Kafka records contract. |

## github.com/faustbrian/go-cloudevents/adapters/outbox/v2

Map CloudEvents explicitly to and from Golib transactional outbox envelopes while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/outbox/v2`](../adapters/outbox/outbox.go) | Use when an application owns both a CloudEvent and Golib transactional outbox envelopes contract. |

## github.com/faustbrian/go-cloudevents/adapters/queue/v2

Map CloudEvents explicitly to and from Golib queue jobs while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/queue/v2`](../adapters/queue/queue.go) | Use when an application owns both a CloudEvent and Golib queue jobs contract. |

## github.com/faustbrian/go-cloudevents/adapters/rabbitstream

Map CloudEvents explicitly to and from Golib RabbitMQ Streams messages while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/rabbitstream`](../adapters/rabbitstream/rabbitstream.go) | Use when an application owns both a CloudEvent and Golib RabbitMQ Streams messages contract. |

## github.com/faustbrian/go-cloudevents/adapters/schema-registry/v3

Map CloudEvents explicitly to and from Golib schema-registry JSON Schema resolution while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/schema-registry/v3`](../adapters/schema-registry/schema_registry.go) | Use when an application owns both a CloudEvent and Golib schema-registry v3 JSON Schema resolution contract. |

## github.com/faustbrian/go-cloudevents/adapters/telemetry

Map CloudEvents explicitly to and from Golib telemetry trace context while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/telemetry`](../adapters/telemetry/telemetry.go) | Use when an application owns both a CloudEvent and Golib telemetry trace context contract. |

## github.com/faustbrian/go-cloudevents/adapters/tenancy/v2

Map CloudEvents explicitly to and from Golib tenant routing identity while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/tenancy/v2`](../adapters/tenancy/tenancy.go) | Use when an application owns both a CloudEvent and Golib tenant routing identity contract. |

## github.com/faustbrian/go-cloudevents/adapters/workflow/v2

Map CloudEvents explicitly to and from Golib durable workflow history while preserving target-owned state and reporting conversion loss.

| Import path | Responsibility |
| --- | --- |
| [`github.com/faustbrian/go-cloudevents/adapters/workflow/v2`](../adapters/workflow/workflow.go) | Use when an application owns both a CloudEvent and Golib durable workflow v2 history contract. |

Choose the next guide in the [documentation index](README.md)
or return to the [repository entry](../README.md) for installation and ownership.
