# Contributing

## Before Editing

1. Read [`AGENTS.md`](AGENTS.md) and the affected module's docs.
2. Run the affected module's deterministic Go tests as described below.
3. Identify owned dependencies and reverse dependants in `modules.json`.
4. Preserve unrelated work and generated/corpus provenance.

## Changes

Keep commits focused and conventional. Update every affected changelog with
the behavior and migration impact. Public API changes require compatibility
evidence and documentation. Specification behavior requires a decision record,
fixture coverage, and interoperability evidence.

New direct dependencies and dependency updates must follow the
[dependency governance policy](AGENTS.md#dependencies-and-supply-chain). Package-local
update bots are forbidden; the root policy owns every module and action update.

Specification-backed changes must follow the
[specification governance contract](AGENTS.md#design), update
the affected [stable decision entries](docs/specification-decisions.md), and
complete the Specification Decisions
section of the pull request template. An unresolved interpretation or stale
source pin is release-blocking; peer behavior cannot silently select policy.

Selected full/release mutation gates must finish with zero surviving viable
mutants; mutation is not required for an ordinary local or pull-request check.

Do not fork shared verification workflows, add permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use the checksum-pinned public Tools v2.4.0 release declared in `.golib.yaml`.
Download its released `golib` binary rather than building the verifier from source.

Use Go 1.27.2 for the producer; module/API minima remain Go 1.27.0.
An isolated checkout needs no sibling repositories or coordination goal files.
Run deterministic tests in the affected module (root shown):

```bash
GOWORK=off go test ./...
GOWORK=off go test . -run '^$' -fuzz '^FuzzDecodeHTTP$' -fuzztime=2s -timeout=2m
```

The first command includes fuzz seeds/corpus, not an active fuzz campaign.
The second explicitly fuzzes one native target. For a target adapter, run Go
tests from that adapter's directory with `GOWORK=off`; public dependencies are
used instead of sibling replacements. Exercise `integration/target-adapters`
when its owned composition changes.

Install the verified public `golib` release declared by `.golib.yaml` for
repository checks; CI verifies the same release without building it from source:

```bash
make check
make ci
golib specification check
```

`make check` and `make ci` run bounded local checks for all catalog modules;
`ci` also checks repository/cohesion structure. Base-free local checks run
deterministic fuzz seeds/corpus, while changed-source CI selects active smoke.
The specification command checks the local decision register without fetching
upstream sources. No module declares required services or test tags.

Optional one-package mutation, when editing the deprecated facade:

```bash
golib mutation --module adapters/golib --package .
```

For selected full qualification and releases, retain the original campaign:

```bash
make check-full
make -f verification/package.mk conformance
make -f verification/package.mk interoperability
golib specification check --online
golib api check
golib release dry-run --module .
```

Full checks retain all eleven native fuzz targets at 100x/2m each, benchmarks,
conformance and interoperability. Conformance fetches the nineteen immutable
official sources and verifies local fixtures; JavaScript interoperability uses
the pinned Node.js/SDK runtime. These are explicit online/full boundaries, not
ordinary development prerequisites. Select the affected adapter directory with
`--module` for its release rehearsal; this can run locally and is not
a CI-only release mechanism. Report unavailable or failing commands and never
describe bounded checks alone as release qualification.

## Selected Monthly Assurance

The owner-level [selected-assurance caller](.github/workflows/assurance.yml) is
configured for main only, on day 17 of each month at 05:07 UTC. It selects one
of the 4 current declared root-module fuzz targets using the shared UTC-week
rotation, with 30 seconds of active exploration. The selected risks are:
bounded JSON event/batch decoding and round trips, plus HTTP/Kafka binding
input handling.
HTTP headers/bodies and Kafka records are in-memory inputs; no broker or
live HTTP service is selected. The deprecated facade's seven fuzz targets,
target adapters and integration composition remain separate, with their
original full/manual qualification.

Base-free ordinary checks execute deterministic seeds, not active exploration;
committed-range CI selects affected active smoke. Manual assurance requires an
exact anchored selector from the root fuzz operation in `.golib.yaml`, for
example `^FuzzDecodeJSONBatch$`.
Empty, unknown and malformed selections refuse. The shared job bounds setup
and execution to ten minutes, active workers and CPU to two, owned memory to
2 GiB, native output to 1 MiB and private-cache lifetime. Failure, cancellation,
exhausted budgets and cleanup refusal remain attributable failures. This lane
does not run mutation or replace Required, original full fuzz budgets,
explicit `make check-full`, or relevant release qualification.

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.

## Source Documentation

Install the checksum-pinned public verifier declared in `.golib.yaml`.
`make source-docs-check` validates objective Go comment, marker, link and
generated-producer structure; `make check`, `make check-full` and `make ci`
require it before their existing checks. Shared CI runs the same strict
contract and includes its result in Required. This is structural enforcement,
not a substitute for source-truth review or behavioral verification.
