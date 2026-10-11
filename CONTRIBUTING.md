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

Do not add package-local workflows, permanent replacements, machine-specific
paths, bypass flags, broad mutation exclusions, or aggregate quality metrics
that hide a failing package.

## Verification

Use the checksum-pinned public Tools v2.1.3 release declared in `.golib.yaml`.
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

## Adding A Module

Follow [repository structure policy](AGENTS.md#repository-structure). New modules
require an explicit purpose, ownership boundary, dependency review, package
catalog entry, full quality gates, documentation, changelog, license, security
policy, compatibility plan, and release dry-run.
