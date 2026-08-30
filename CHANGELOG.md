# Changelog

All notable changes follow Keep a Changelog. This module uses semantic
versioning.

## Unreleased

### Changed

- Replace the repository-local verification implementation with the pinned
  `go-library-tools` v1.0.13 CLI and reusable workflow while preserving module
  policy, interoperability fixtures, and content-addressed evidence.
- Use canonical public module checksums for the nested Golib adapter instead of
  bootstrap-only archives.
- Pin CI specification enforcement to canonical shared tooling commit
  `2dd7ae5aa634c99dc9aa9033d361019e5ffc9988`.

### Documentation

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
