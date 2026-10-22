# Changelog

All notable changes follow Keep a Changelog. This module uses semantic
versioning.

## 1.1.3 - 2026-10-22

### Changed

- Publish the current source and documentation with a fresh immutable tag.
- Build with Go 1.27.3 to address GO-2026-6809 in `crypto/x509`.
  Consumers should use Go 1.27.3 or newer.
- Publish independently versioned modules from the same source snapshot:
  - `adapters/golib/v4.1.0`.
  - `adapters/audit/v2.1.0`.
  - `adapters/correlation/v1.1.0`.
  - `adapters/event-sourcing/v2.1.0`.
  - `adapters/jsonschema/v2.1.0`.
  - `adapters/kafka/v1.1.0`.
  - `adapters/outbox/v2.1.0`.
  - `adapters/queue/v2.1.0`.
  - `adapters/rabbitstream/v1.1.0`.
  - `adapters/schema-registry/v3.1.0`.
  - `adapters/telemetry/v1.1.0`.
  - `adapters/tenancy/v2.1.0`.
  - `adapters/workflow/v2.1.0`.

## 1.1.2 - 2026-10-13

This entry records a prepared root patch; publication is separate. Adopt
`github.com/faustbrian/go-cloudevents@v1.1.2` only after its public tag and module
are available. Independently versioned adapters need an explicit corrected root
selection in the consuming application; their existing releases are not republished.

### Fixed

- Restore canonical received base64 admission across JSON and structured
  transports: reject line breaks and null binary members while preserving
  valid JSON escapes and empty binary strings. Reject null JSON batches;
  empty arrays remain valid. Invalid input returns no partial event or batch.

### Source accountability

This root source includes the following 35 non-merge repository commits after
public root v1.1.1 (`5559c521abeb0ac979a06a9405185b8486033c88`) through
`925df9881d1fe33e510c802efda0aa4ff90f4141`. Several concern independently released
adapter families or tooling, not a new root runtime feature. The current
admission repair is outside this frozen range; other pending Unreleased entries
and prior release notes are retained.

- [141a6536f2](https://github.com/faustbrian/go-cloudevents/commit/141a6536f2cbf91836c21c965c08e86e806ed1b1),
  [3a1d8fb457](https://github.com/faustbrian/go-cloudevents/commit/3a1d8fb457b59cc50b29284a4e4a93aa60f898ac),
  [4b5ac41915](https://github.com/faustbrian/go-cloudevents/commit/4b5ac419159624e1550d35ab50aca5bb1e145c7b).
- [fa1bb5fb9e](https://github.com/faustbrian/go-cloudevents/commit/fa1bb5fb9e883e72b05074dbddc4250c45dac7a3),
  [8aa80ae9c4](https://github.com/faustbrian/go-cloudevents/commit/8aa80ae9c4409053aee971af0e67f67d95cf7875),
  [91bac0bd23](https://github.com/faustbrian/go-cloudevents/commit/91bac0bd2348e0006d4db8b1239babc2e4b6c346).
- [0aaa6e6a24](https://github.com/faustbrian/go-cloudevents/commit/0aaa6e6a24c698a11cf9dd6d0bdeb7363a3be8dd),
  [86a8d93a21](https://github.com/faustbrian/go-cloudevents/commit/86a8d93a213f26d9fa7008afe2713d571801fa9f),
  [17e1234c01](https://github.com/faustbrian/go-cloudevents/commit/17e1234c015ecec77b40d36418d5dfa5f18b479e).
- [6cb28d8ef9](https://github.com/faustbrian/go-cloudevents/commit/6cb28d8ef9e4ec989e28d0897336d331d9af3562),
  [6328dd1726](https://github.com/faustbrian/go-cloudevents/commit/6328dd17262d10890b9c3f2c781602060a3354bf),
  [f7e18aee23](https://github.com/faustbrian/go-cloudevents/commit/f7e18aee23c5cc59f5c638067fe92737429440c1).
- [c35e7e2149](https://github.com/faustbrian/go-cloudevents/commit/c35e7e214908459c2b80467ea6be421f1d537b26),
  [bb458dba14](https://github.com/faustbrian/go-cloudevents/commit/bb458dba144ad09a7a62ef3d78a48ea95ccb4bde),
  [9a4a405070](https://github.com/faustbrian/go-cloudevents/commit/9a4a405070a24f4eea65856c9a5c9c6543029af3).
- [ad718ca37e](https://github.com/faustbrian/go-cloudevents/commit/ad718ca37e77d64e73716a5f2ac1afc5c4aaa85c),
  [55dc285d2d](https://github.com/faustbrian/go-cloudevents/commit/55dc285d2d1d368c6a95bb8277add683b82cc25a),
  [68df4d8273](https://github.com/faustbrian/go-cloudevents/commit/68df4d827377d487d72c42d4611b5ee0f9de4da9).
- [32c81411ab](https://github.com/faustbrian/go-cloudevents/commit/32c81411ab3ec521b48558c158179da655a11f20),
  [f7165a0ec5](https://github.com/faustbrian/go-cloudevents/commit/f7165a0ec56b88bb0d6975716d29ebe6e88c8749),
  [9fbe050c73](https://github.com/faustbrian/go-cloudevents/commit/9fbe050c73e8d19600f922a8952230deed35e68e).
- [ce4d6f3deb](https://github.com/faustbrian/go-cloudevents/commit/ce4d6f3deb6bf3221ac6db6b4c7a8c0301b10615),
  [2f5217bdfa](https://github.com/faustbrian/go-cloudevents/commit/2f5217bdfa060453f19e62ac8fc76ee245c86b0a),
  [c81c0870ee](https://github.com/faustbrian/go-cloudevents/commit/c81c0870ee92e6cb9b069a09848f928202077ee9).
- [c463941dd2](https://github.com/faustbrian/go-cloudevents/commit/c463941dd20825a67d94d73e51d61ce47432eddd),
  [0a1c18fd37](https://github.com/faustbrian/go-cloudevents/commit/0a1c18fd37029dd1cc133a0dd1eb9ad76698eb97),
  [87d5db0c9e](https://github.com/faustbrian/go-cloudevents/commit/87d5db0c9ee37010ecfce8cfce7dda7f93c0398e).
- [6256325388](https://github.com/faustbrian/go-cloudevents/commit/6256325388ead6dc2a99dc2850a0c5b375860134),
  [bacfb75212](https://github.com/faustbrian/go-cloudevents/commit/bacfb752125817076fbfeb951879d30474199b43),
  [f8755fef42](https://github.com/faustbrian/go-cloudevents/commit/f8755fef42b6f629a670e347ffdbaf43a79c70f8).
- [faa16c8a08](https://github.com/faustbrian/go-cloudevents/commit/faa16c8a08c39cb237ad3c27fbc826537fd57573),
  [7d8554ce20](https://github.com/faustbrian/go-cloudevents/commit/7d8554ce20cf332012fd110d217c700be15c825a),
  [9df39dd95a](https://github.com/faustbrian/go-cloudevents/commit/9df39dd95a15a3dd8f143e9a2f58aa1f9ca7a1b1).
- [bc98b48945](https://github.com/faustbrian/go-cloudevents/commit/bc98b48945e61da4cb691093a03a77fa944b1dc7),
  [925df9881d](https://github.com/faustbrian/go-cloudevents/commit/925df9881d1fe33e510c802efda0aa4ff90f4141).

## Unreleased

### Changed

- Select the patched documentation spelling parser in the repository lock
  and immutable shared CI tooling. Build the exact tooling source with public
  dependency verification; published Go module dependencies are unchanged.

- Align reusable CI and its executed setup tooling to one immutable source,
  retaining the configured v1.8.4 CLI and independent security qualification.

- Adopt published Audit2, Queue2, EventSourcing2 and Tenancy2 adapters in the
  target integration recipes. The deprecated facade-v3 identity also adopts
  public Schema Registry2 nominal configuration while retaining mapping and
  ownership semantics; prior facade releases remain available.

- Declare current Audit v2, Queue v2 and Tenancy v2 adapter identities
  consistently in the family catalog while preserving legacy facade inputs.
- Add isolated native family security qualification without substituting for
  ordinary required CI or release checks.

- Upgrade checksum-pinned `go-library-tools` and reusable CI to v1.8.4 so
  selected adapter releases follow the proportional assurance policy while
  retaining every configured verification gate.

### Documentation

- Clarify constructor/encoder versus decoder admission, copied-value privacy,
  opaque payload/extension loss, cooperative readers/validators and exact
  transport-state ownership without changing wire acceptance or source pins.

- Record versioned family trust boundaries and accepted collaborator risks
  with explicit ownership, mitigations and review conditions.

- Document the public Workflow v2 adapter and explicit facade-v2 and
  non-releasable integration adoption without changing the CloudEvents root
  module contract or unrelated dependency selections.

- Document staged adoption of the independently versioned Workflow v2 adapter;
  the CloudEvents root and released-v1 facade and integration dependencies
  remain unchanged.
- Require Go 1.27.0 across the repository's module language, minimum
  compatibility, development, and CI toolchain claims.
- Mark the broad `adapters/golib` compatibility facade as deprecated for new
  adoption and point applications to the target-oriented adapters and their
  non-releasable Kafka and schema validation recipe.

## 1.1.0 - 2026-09-09

### Added

- Add canonical optional-adapter error categories and explicit conversion-loss
  value types so independently released target adapters preserve one public
  classification contract without a shared runtime.
- Add independently adoptable target adapter modules for audit, correlation,
  event sourcing, JSON Schema, Kafka, outbox, queue, RabbitMQ Streams, schema
  registry, telemetry, tenancy, and workflow boundaries.

### Changed

- Upgrade checksum-pinned `go-library-tools` and reusable CI to v1.6.2 so
  valid nested-module documentation links remain supported.

- Upgrade checksum-pinned `go-library-tools` and reusable CI to v1.4.0 while
  retaining schema-v2 cohesion metadata and frozen specification decisions.

- Adopt schema-v2 cohesion metadata for both releasable modules, the pinned
  `go-library-tools` v1.3.0 CLI, repository-local cohesion validation, and the
  immutable shared workflow that enforces it in CI.
- Adopt the pinned `go-library-tools` v1.2.0 CLI and reusable workflow so CI
  enforces specification decisions, conformance bindings, source monitoring,
  and change control while preserving module-owned policy, interoperability
  fixtures, and content-addressed evidence.
- Use canonical public module checksums for the nested Golib adapter instead of
  bootstrap-only archives.

### Documentation

- Bind both releasable modules' ecosystem and protocols-and-descriptions family
  navigation to the immutable v1.4.0 documentation.
- Link both releasable modules to the versioned Golib ecosystem index and
  shared package-selection, construction, ownership, and lifecycle guidance.
- Replace archived monorepo links with a package-owned documentation index.
- Enforce the [specification decision register](docs/specification-decisions.md),
  conformance bindings, immutable authority pins, release monitoring, and
  append-only decision history through the shared repository check.
- Current decision records: `CLOUDEVENTS-DEC-001 sha256:317e42da14b62ad0511bc935463c8fb01cf037db77e0c20e9df885ea01559b7a`;
  `CLOUDEVENTS-DEC-002 sha256:6c35a0124b6ca9a733f2cd2c8e2c4431997b3f3b6527193232d596ba2c177bdb`;
  `CLOUDEVENTS-DEC-003 sha256:695e83ba9df4c8eb1f986ea49b7194934a894f2a3b77503a8c7fe1bfa533de88`;
  `CLOUDEVENTS-DEC-004 sha256:bb3e9e20e33e1619a51406e2c30fbb1db462138aeb005e9ff06904de925a5345`;
  `CLOUDEVENTS-DEC-005 sha256:c15f47585a8a86e8382595cceece5e71d5302b963bc7c35a36a5df501b572415`;
  `CLOUDEVENTS-DEC-006 sha256:7f3898d219e264bcb7fb23cf3493e4ecf51bd4d2512d38cf1d3ae4d36308e009`;
  `CLOUDEVENTS-DEC-007 sha256:dd58a6179aa2c6674f8f75b8bd0792667548d9a288e346800c633a76826ae38e`;
  `CLOUDEVENTS-DEC-008 sha256:58bcf8a8b3ddf3d11477544af9b1e418a8605f3170dbcbc372e7394ef7323c03`;
  `CLOUDEVENTS-DEC-009 sha256:c1545df6d6f51dd4f93d7dca7dfae4964b7ee7a0650cdc4bc2036f791fdcedb2`;
  `CLOUDEVENTS-DEC-010 sha256:7a33d17a4c654d30849ead7ddea51839e2c006f8129e2e67781b9d389525ca7b`;
  `CLOUDEVENTS-DEC-011 sha256:3bdde5243c68f63da55fdfc9d00f2a693b15d81461349cad37cab527c7ed66b6`;
  `CLOUDEVENTS-DEC-012 sha256:48d7db91a9d4bc2ca4d7858b20ee6ebb5978960322662faf590445e5a85502fe`.

## 1.0.0 - 2026-08-25

### Changed

- Refresh the pinned JavaScript lockfile evidence and standalone adapter
  checksum after the security-fixed interoperability dependency update.

- Upgrade Go cryptography and network dependencies and pin JavaScript
  interoperability to the security-fixed `uuid` release.

- Exclude intentional nested modules from root local-proxy archives so local,
  bootstrap, CI, and public module checksums describe the same source
  boundary.

- Track the pinned documentation-tool lockfile so clean CI checkouts install
  the exact validated cspell dependency.

- Reconcile standalone dependency checksums against deterministic current
  module archives so CI, local verification, and release consumers resolve
  identical content.

- Harden standalone documentation validation with deterministic spelling and
  link checks, package-specific documentation gates, and repository-local
  contributor guidance.

### Documentation

- Link the package README to package-owned documentation.

### Added

- Canonical specification decision register covering stable-line errata, data
  presence, extension typing, duplicate metadata, deterministic JSON, HTTP and
  Kafka binding conflicts, resource ownership, and explicit schema validation.
- Immutable CloudEvents 1.0 event, data, and typed context-attribute model.
- Deterministic JSON event and batch encoding with bounded hostile-input
  decoding.
- Transport-neutral HTTP and Kafka binary and structured content-mode mappings.
- Selected distributed-tracing and partitioning extension validation.
- Explicit caller-supplied schema validation without implicit I/O.
- Pinned normative matrix, errata decisions, provenance, fuzzing, stress,
  ownership, interoperability, and benchmark coverage.
- Bidirectional Go and JavaScript SDK interoperability fixtures for JSON,
  batch, HTTP, Kafka, tracing, partitioning, and unknown extensions.
- Strict and loss-aware JSON, batch, HTTP, and Kafka encoders with deterministic
  reports for metadata materialization, abstract extension-type normalization,
  and unrepresentable JSON payload whitespace.
- Requirement-level normative and explicit unsupported-surface matrices,
  archived official-kit report, and checksum-pinned task-owned Node.js runtime
  for independent JavaScript interoperability.

### Changed

- Publish the module from its standalone `github.com/faustbrian/go-cloudevents` identity while preserving its documented API and behavior.
- Preserve declared JSON payload bytes across structured JSON, HTTP, Kafka, and
  batch round trips instead of compacting payloads during encoding.
- Assert complete JavaScript SDK consumer context, extension, and semantic
  payload results, including its explicit HTTP timestamp, null, empty-data, and
  default-content-type normalizations.
- Reject or ignore hostile HTTP headers without lowercasing unbounded unowned
  names; CloudEvents attribute-name limits are enforced before case folding.
- Treat `text/*` binary HTTP and Kafka payloads as text data, and require an
  explicit loss report whenever encoding must materialize `application/json`,
  `text/plain`, or `application/octet-stream` to preserve a runtime data kind.
- Reject explicit content-type and data-kind conflicts from strict JSON, HTTP,
  and Kafka encoders instead of allowing a later decoder to reinterpret or
  reject the payload.
- Compare official Go SDK decoding against the same canonical byte corpus and
  retain provider-specific interoperability conflicts as explicit results,
  including `SetData(nil)` normalization to absent data.

### Security

- Bound event bytes, data bytes, attribute counts and sizes, JSON depth, and
  batch size before retaining untrusted data.
- Bound HTTP context-attribute names and all Kafka record metadata copied by
  the decoder.
- Reject duplicate and conflicting metadata, invalid Unicode, malformed URI
  values, invalid media types, and non-canonical base64.
- Reject non-ASCII distributed-tracing `tracestate` values.
- Apply the configured attribute-value limit to unknown JSON extensions before
  retaining decoded metadata.
- Bound HTTP `Content-Type` before media-type parsing and reject excess or
  duplicate binding metadata before retaining decoded attributes.
