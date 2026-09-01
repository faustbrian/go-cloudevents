# Specification conformance decisions

The [decision register](../docs/specification-decisions.md) owns observable
interpretations. The [normative requirement matrix](normative-requirements.json)
owns statement-level coverage, while this index binds every stable decision to
the executable conformance manifests without claiming external certification.

- `CLOUDEVENTS-DEC-001` — stable-line clarifications after 1.0.2
- `CLOUDEVENTS-DEC-002` — extension attribute-name compatibility
- `CLOUDEVENTS-DEC-003` — data presence and JSON representation
- `CLOUDEVENTS-DEC-004` — unknown extension abstract types
- `CLOUDEVENTS-DEC-005` — duplicate and null context attributes
- `CLOUDEVENTS-DEC-006` — URI and URI-reference input strictness
- `CLOUDEVENTS-DEC-007` — deterministic JSON bytes
- `CLOUDEVENTS-DEC-008` — structured metadata conflicts and mode selection
- `CLOUDEVENTS-DEC-009` — HTTP body presence, ownership, and cancellation
- `CLOUDEVENTS-DEC-010` — Kafka tombstones, keys, and batch scope
- `CLOUDEVENTS-DEC-011` — explicit schema validation without implicit I/O
- `CLOUDEVENTS-DEC-012` — distributed tracing lifecycle ownership

The maintained Go and JavaScript SDK comparisons are differential evidence,
not normative authorities. Official conformance features remain a separate
fixture-evidence lane. Source pins and mutable release monitoring remain
separate in [`monitoring.json`](monitoring.json).
