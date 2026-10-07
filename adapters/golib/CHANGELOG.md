# Changelog

All notable changes to this module are documented here.

## Unreleased

### Changed

- Adopt Correlation v1.1.2 and its indirect Identifier v2 dependency
  while retaining facade v3 correlation and audit metadata contracts.

- Select Kafka root module v1.1.0 while preserving record ownership,
  transport metadata, and CloudEvents binding behavior.

### Changed

- Adopt Queue v1.1.2 in the facade dependency graph while retaining the
  public facade v3 conversion contracts.

- Adopt RabbitMQ Streams v1.1.1 and record Snappy v1.0.0 in the facade
  dependency graph while retaining the public facade v3 identity.

- Adopt the public Tenancy v2, EventSourcing v2 and Schema Registry v2 contracts
  through facade `github.com/faustbrian/go-cloudevents/adapters/golib/v3`.
  Update facade and corresponding core imports together; tenant values, audit
  metadata, event-store messages/state and registry configuration now use their
  actual public v2 nominal identities. Conversion, trust, retained ownership and
  loss semantics remain unchanged. Published facade v1 and v2 remain available.

- Adopt the actual public Workflow v2 adapter and core through the independent
  `github.com/faustbrian/go-cloudevents/adapters/golib/v2` module. Update facade
  and Workflow imports together: history inputs, outputs and retained definition
  references now use Workflow v2 nominal types. Conversion ownership, payload
  presence, loss reporting and unrelated direct dependency selections remain
  unchanged. Released facade v1 stays available; new code should select only
  the target-oriented adapters it needs.
- Enable the existing facade mutation gate for the independent v2 release
  policy.

### Deprecated

- Deprecate the broad compatibility facade for new adoption while retaining it
  throughout v1. Migrate to the target-oriented adapters in the parent adoption
  guide; the facade is excluded from the recommended set because it owns 26
  module dependencies.

### Documentation

- Relocate the Kafka and schema validation recipe to the non-releasable
  target-adapter integration module so recommended composition no longer
  imports the deprecated facade.

## 1.1.0 - 2026-09-09

### Changed

- Preserve the released umbrella import path as a compatibility facade while
  delegating behavior to independently released target adapter modules.
- Require CloudEvents v1.1.0 and target adapter v1.0.0 module identities.
- Upgrade checksum-pinned repository tooling and reusable CI to v1.4.0 while
  retaining the adapter's schema-v2 cohesion and specification contracts.
- Replace bootstrap-only archive checksums for all owned dependencies with
  their canonical public v1.0.0 module identities for clean consumer
  verification.

### Documentation

- Add an executable Kafka, schema-registry, JSON Schema, and CloudEvents
  composition with explicit validation, settlement, failure, and shutdown
  ownership.
- Bind ecosystem and protocols-and-descriptions family navigation to the
  immutable v1.4.0 documentation.
- Publish schema-v2 family, selection, ownership, lifecycle, compatibility,
  and documentation metadata and link to the versioned Golib ecosystem index.
- Move detailed module guidance behind a concise README and documentation index.
- Use human-oriented section names and package-owned documentation links.

## 1.0.0 - 2026-08-25

### Documentation

- Link the package README to package-owned documentation.

### Added

- Add explicit, loss-reporting Golib conversions for event sourcing, outbox,
  queue, workflow, Kafka, correlation, tenancy, telemetry, and audit metadata.
- Add caller-invoked JSON Schema and schema-registry validation with static URI
  mappings, bounded cache resolution, explicit availability policy, and no
  event-controlled lookup target.
- Document collision, trust, ownership, loss, migration, and security policy.
- Enforce Kafka record limits before copying consumed key, value, or header
  metadata.
- Add a structured JSON RabbitMQ Streams mapping that bounds copied transport
  metadata and keeps broker state and properties outside the CloudEvents
  context.
- Require tenant metadata during trusted tenant extraction and validate audit,
  queue, and event-sourcing tenant identifiers before emission or replay
  acceptance; tenant remains optional for audit records that do not declare it.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-cloudevents/adapters/golib` identity while preserving its documented API and behavior.
- Refresh local `v0.0.0` owned-module checksums after dependency manifests and
  release notes were normalized; runtime behavior and public APIs are
  unchanged.
- Align the transitive `golang.org/x/text` dependency with v0.41.0 after the
  owned module graph removed GO-2026-5970.
- Replace direct `RegistryJSONSchemaValidator` field initialization and the
  `LookupSchema` callback with `NewRegistryJSONSchemaValidator` and
  `RegistryJSONSchemaConfig`. Callers must provide a bounded `ResolveCache`, a
  static URI-to-lookup map, an explicit availability policy, and a positive
  timeout. This is a source-breaking pre-v1 migration that prevents mutable or
  event-controlled resolver selection.
- Sort extension conversion losses by field for deterministic reports.
- Report payload-kind, content-type, and schema loss across canonical envelope
  conversions, and enumerate every unselected canonical audit field.
- Reject retained queue content-type collisions instead of silently replacing
  canonical queue metadata.
- Reject event-sourcing subjects and event-sourcing or queue extensions that
  disappear while their non-empty retained canonical values remain available.
- Preserve nil and present-empty payloads distinctly across outbox, queue, and
  workflow round trips; workflow callers must retain `DataWasNil` with the
  other `WorkflowState` fields.
- Reject non-canonical workflow event type suffixes instead of silently
  normalizing them during a round trip.
- Allow trusted audit metadata extraction to preserve an absent optional tenant.

### Distribution

- Include the canonical MIT licence for independent publication.

### Fixed

- Keep RabbitMQ Streams decoding source formatted identically by the supported
  Go 1.26.6 toolchain and Go 1.27 without changing decoded values or state.
- Return the CloudEvents context-required category from both schema validators
  when either is called directly with a nil context, avoiding a registry panic.
- Preserve the tenant-specific error classification when malformed queue
  tenant metadata also violates generic queue metadata validation.
