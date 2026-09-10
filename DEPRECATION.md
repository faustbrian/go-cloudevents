# Deprecation Policy

Deprecations MUST identify the replacement, reason, migration steps, and
earliest removal version. Public Go identifiers use a valid `Deprecated:` doc
paragraph and corresponding changelog entry.

At `v1` and later, a supported replacement SHOULD exist for at least one minor
release before removal. Security or correctness defects MAY require faster
removal when continued support would be unsafe; the release notes must explain
the exception.

Silent behavior changes, undocumented aliases, and indefinite deprecated code
are prohibited. Deprecations are checked during compatibility and release
review.

## `adapters/golib` compatibility facade

`github.com/faustbrian/go-cloudevents/adapters/golib` is deprecated for new
adoption. It is retained for v1 compatibility, but its broad integration
surface owns 26 module dependencies and is therefore excluded from the
recommended adapter set.

Replace each facade import with only the target-oriented adapters required by
the application boundary. The root [adoption guide](README.md#adoption-guidance)
lists the supported replacements for audit, correlation, event sourcing, JSON
Schema, Kafka, outbox, queue, RabbitMQ streams, schema registry, telemetry,
tenancy, and workflow integrations. Migrate one boundary at a time; the target
adapters preserve the same public mapping and validation contracts without
requiring the unrelated bridge dependencies.

The earliest possible removal is `adapters/golib/v2.0.0`. No removal is
scheduled, and the facade remains supported for the v1 compatibility line.
