# Engineering Policy

The key words "MUST", "MUST NOT", "REQUIRED", "SHALL", "SHALL NOT",
"SHOULD", "SHOULD NOT", "RECOMMENDED", "NOT RECOMMENDED", "MAY", and
"OPTIONAL" in this document are to be interpreted as described in BCP 14
[RFC2119] [RFC8174] when, and only when, they appear in all capitals, as
shown here.

## Scope And Authority

- This file is the canonical policy for the complete repository.
- Package policies MAY add stricter domain rules but MUST NOT weaken this file.
- `CLAUDE.md` and tool-specific files MUST point here rather than duplicate it.
- Historical `.ai/GOAL*.md` files are requirements and evidence, not proof of
  completion. Current executable evidence is REQUIRED.

## Repository Structure

- The public root module MUST live at the repository root.
- Intentional optional or test modules MAY live in explicit nested directories.
- Commands MUST live under `cmd/`; private shared code MUST live under
  `internal/`; root automation MUST live under `scripts/`.
- Public module paths MUST match their repository-relative directories beneath
  the module path declared by the root `go.mod`.
- Every module MUST be declared in `modules.json`, and every package MUST be
  declared in `packages.json`.
- Independently releasable modules MUST retain independent `go.mod` files and
  directory-prefixed semantic-version tags.
- Cross-module dependencies MUST remain acyclic and MUST use public contracts.
- Permanent `replace` directives, sibling repositories, and absolute developer
  paths are forbidden in releasable modules.

## Design

- Prefer standard-library interfaces and explicit composition over hidden
  registration, global state, reflection-driven wiring, or service locators.
- Public APIs MUST make ownership, cancellation, retries, timeouts, resource
  limits, error semantics, and concurrency behavior observable.
- Interfaces SHOULD be defined by consumers and MUST remain narrowly scoped.
- Optional integrations SHOULD be adapters or nested modules, not mandatory
  dependencies of a core package.
- Breaking protocol or specification ambiguities MUST be documented as explicit
  decisions and covered by tests.

## Safety And Concurrency

- Shared mutable state MUST have one documented synchronization owner.
- Goroutines MUST have explicit lifetime, cancellation, and shutdown.
  Goroutine lifecycle changes MUST include targeted leak tests.
  Fire-and-forget goroutines are forbidden.
- Channels MUST have documented ownership and closure rules.
- Locks MUST NOT be held across caller callbacks, network IO, blocking channel
  operations, or unbounded work.
- Every external operation MUST accept or derive a bounded `context.Context`.
- Response bodies, files, rows, transactions, timers, tickers, connections,
  and temporary resources MUST be closed on every path.
- Integer conversions, sizes, offsets, recursion, decompression, and allocation
  from untrusted input MUST be bounded before allocation or conversion.
- Secrets and credentials MUST NOT appear in errors, logs, traces, snapshots,
  fixtures, mutation reports, or generated artifacts.

## Testing

- Behavioral changes MUST include meaningful tests before completion.
- Tests MUST assert outcomes, invariants, errors, cleanup, and state transitions;
  line execution without behavioral assertions is not acceptable coverage.
- Verification MUST follow the proportional assurance tier of the actual
  change:
  - **Tier A** documentation, metadata, registration, and generated
    documentation changes require affected structural validation and diff
    review.
  - **Tier B** internal behavioral changes require focused behavior tests,
    affected package or module tests, applicable formatting and static checks,
    and one complete review.
  - **Tier C** public API, lifecycle, security, persistence, and concurrency
    changes require regression or characterization tests, focused behavior,
    API compatibility where applicable, affected package and integration
    tests, direct owned reverse consumers, and one independent complete-diff
    review.
  - **Tier D** releases and ecosystem milestones require only the relevant
    compatibility, composition, consumer, and aggregate checks over immutable
    release inputs.
- Coverage MUST be sufficient to prove the affected observable behavior.
  Exact statement coverage is required only when an applicable gate explicitly
  selects it for a material risk.
- Mutation, fuzz, race, leak, performance, conformance, external-service,
  clean-consumer, release-rehearsal, and aggregate checks MUST run only when
  they exercise a named material risk or the applicable Tier D boundary.
- Concurrent code changes MUST pass the race detector and targeted stress or
  leak tests when those checks exercise the changed concurrency boundary.
- Parser or hostile-boundary changes MUST include resource limits and
  deterministic regressions; fuzzing is required only when it materially
  exercises the changed parser risk.
- Specification claims MUST be proven against pinned official fixtures and
  independent implementations where applicable.
- Performance claims MUST be supported by benchmarks that compare equivalent
  behavior and record enough environment and corpus detail to reproduce them.

## Required Commands

- `make inventory` validates repository and package manifests.
- `make check` runs the selected repository contract when its gates match the
  change's assurance tier.
- `make ci` runs the configured repository contract; it MUST NOT be treated as
  a universal requirement for unrelated or unchanged modules.
- Local commands and CI MUST use the same scripts and thresholds.
- Missing tools, services, packages, profiles, or reports MUST fail only when
  the selected applicable gate requires them.
- NilAway is advisory; its findings MUST remain visible and tracked against a
  no-regression baseline.

## Evidence Validity And Reuse

- Evidence validity MUST be determined by the behavior-affecting inputs of the
  applicable gate, not by a commit hash, branch name, timestamp, or unrelated
  repository-history shape alone.
- Immutable evidence SHOULD be reused when the relevant source, tests,
  dependencies, fixtures, and gate implementation are unchanged.
- Mutable plans, progress notes, reviews, and prose MUST NOT require recursive
  hashes or cryptographic binding.
- Commit hashes MAY be recorded for traceability, but MUST NOT be the sole
  evidence cache key or invalidation condition.
- Go toolchain revision metadata MAY be retained when a build, test, profile,
  or diagnostic artifact requires it. That metadata is descriptive only and
  MUST NOT make an otherwise identical gate-input fingerprint stale.
- A history rewrite, rebase, squash, reset, repository reinitialization,
  metadata-only commit, or unrelated-file change MUST NOT invalidate evidence
  when the complete gate-input fingerprint is unchanged.
- Agents MUST NOT rerun an expensive gate solely to attach an already proven
  result to a new `HEAD`.
- Agents MUST NOT restart the complete package matrix after a force-push,
  rebase, squash, reset, or other history-only change. Previously verified
  package results SHOULD be reused when their relevant inputs are unchanged.
- After a change, agents MUST rerun only the gates, modules, packages, and
  reverse dependants affected by that change.
- Reused evidence MUST retain its original result and MUST NOT be rewritten to
  imply that a gate executed again.
- Long-running applicable checks SHOULD checkpoint independently valid module
  or package results so interruption does not force unrelated reruns.
- Evidence MUST NOT be reused when the relevant input identity is missing or
  ambiguous.
- If repository tooling invalidates evidence solely because `HEAD` changed,
  agents MUST correct the evidence model instead of launching a repository-wide
  rerun with no changed gate inputs.

## CI And Workflows

- `.github/workflows/ci.yml` is the only owned GitHub Actions workflow.
- Package-local workflows MUST NOT be added.
- Actions and external tools MUST be pinned to immutable versions.
- Every module selected by the applicable assurance tier MUST have an
  attributable result.
- The stable required job MUST fail for failed, cancelled, skipped, or missing
  module results.
- Required checks MUST NOT use `continue-on-error`, `|| true`, permissive
  thresholds, or warning substitutions.

## Dependencies And Supply Chain

- Dependencies MUST be necessary, maintained, license-compatible, and pinned to
  reviewed current versions.
- Standard-library functionality MUST NOT be wrapped merely to create an owned
  abstraction; wrappers require a stable policy or portability boundary.
- Generated code and vendored corpora MUST record source, version, checksum,
  license, generation command, and update procedure.
- Release verification MUST select vulnerability, secret, license, SBOM,
  provenance, and clean-consumer checks only when they exercise the release's
  dependency, artifact, or adoption boundary.

## Documentation

- Public identifiers MUST have useful Go documentation describing semantics,
  invariants, ownership, errors, concurrency, and caveats where relevant.
- Comments MUST explain why a constraint or non-obvious implementation exists;
  they MUST NOT narrate obvious syntax.
- Public modules MUST document the entry points, ownership, compatibility,
  adoption, and security information needed for their actual surface.
- Documentation and examples affected by a change MUST be validated; an
  unchanged documentation fleet MUST NOT be revalidated for an unrelated
  change.

## Changelogs

- Material user-visible behavior, API, dependency, compatibility, security, or
  documentation changes MUST update the affected module `CHANGELOG.md` in the
  same commit.
- Entries MUST describe behavior and migration impact, not internal activity.
- Changes to multiple modules MUST update each materially affected changelog.
- Unreleased entries MUST NOT be silently rewritten or removed.
- Generated, dependency, security, compatibility, and deprecation changes are
  user-visible and require entries when they affect consumers or operators.

## Completion

- Run the narrowest affected gates during development and the release gates
  selected by the applicable Tier D risks before declaring completion.
- Re-run affected gates after the final source, test, dependency, documentation,
  workflow, or generated-file change.
- Report exact commands and results. A skipped, blocked, stale, or warning-only
  gate is not a pass.
