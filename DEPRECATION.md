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

All generations of `github.com/faustbrian/go-cloudevents/adapters/golib`
are deprecated for new adoption. Released v1, v2 and v3 retain their original
contracts. The pending facade v4 adopts JSONSchema2, Outbox2 and Registry3
alongside the existing Workflow2, Tenancy2 and EventSourcing2 contracts.
The broad integration
surface owns 26 module dependencies and is therefore excluded from the
recommended adapter set.

Replace each facade import with only the target-oriented adapters required by
the application boundary. The root [adoption guide](README.md#adoption-guidance)
lists the supported replacements for audit, correlation, event sourcing, JSON
Schema, Kafka, outbox, queue, RabbitMQ streams, schema registry, telemetry,
tenancy, and workflow integrations. Migrate one boundary at a time; the target
adapters preserve the same public mapping and validation contracts without
requiring the unrelated bridge dependencies.

No removal is scheduled. Any future removal requires a later major release;
the facade remains present in the prepared `adapters/golib/v4.0.0` module.
