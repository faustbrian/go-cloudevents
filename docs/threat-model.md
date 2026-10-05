# CloudEvents family threat model

Model version: 1.0. Review date: 2026-10-05.
Owner: Brian, repository maintainer. Next review: 2027-01-05, or sooner when
a listed review condition occurs. This model applies to current main's core,
declared adapters and target-adapter integration module. It does not certify
historical release versions or downstream services.

## Assets and actors

Assets are event payloads, tenant and operational identifiers, transport metadata,
schema references, caller contexts, retained envelope state, signing identity,
dependency inputs and CI/release integrity. Event producers and transport peers
can control event bytes, attributes, extension names/values, and copied metadata.
Caller-provided readers, validators, caches, resolvers and propagation policies
are explicit application collaborators, not sandboxed plugins.

## Trust boundaries and controls

| Boundary | Responsibility and controls |
| --- | --- |
| Core JSON, batch, HTTP and Kafka admission | Explicit size/depth/count policies bound retained input. Parsing establishes syntax, never identity or authority. Rejected values are absent from core validation diagnostics. |
| Canonical event and conversions | Owned copies prevent input aliasing. Returned payloads and retained queue/outbox/event-store/workflow state remain sensitive application data. Loss reports identify fields rather than lost values. |
| Tenant, correlation, audit and telemetry extraction | An explicit trusted path controls adoption. Authentication, tenant/resource authorization and propagation policy remain caller-owned; assertions alone never authorize a principal. |
| Schema and registry adapters | Decoding performs no resolution. Explicit validation requires caller-owned resolver/cache policy, deadlines and configured URI mapping. Callers must apply network/redirect/DNS allowlists and response bounds in their collaborators. |
| Queue, outbox, event-sourcing and workflow adapters | Conversion preserves the source envelope's state; it does not enqueue, persist, retry, deduplicate or authenticate. Source packages and application dispatchers own replay/order/dead-letter policy. |
| HTTP/Kafka and transport adapters | No implicit broker or registry I/O. HTTP body interruption depends on the reader; transport credentials and server/client resource limits remain caller-owned. |
| Observability | The library emits no logs/spans/metrics. Applications must use finite codes/categories and avoid payloads, identifiers, raw metadata and delegated error text by default. |
| Build and publication | Pinned actions/tool inputs, normal protected delivery, signed release evidence and public consumer checks establish separate source/artifact boundaries. Maintainer-key compromise is not prevented by a signature. |

Detailed data classification, emission rules and named executable checks are in
[the security review](security-review.md). Limits and cancellation obligations
are in [the security policy](../SECURITY.md).

## Accepted collaborator risks

These are explicit medium residual risks at application boundaries, not evidence
that a failing library gate or known library defect has been waived. Brian owns
maintainer review; each integrating application owns enforcement of its controls.

| Risk | Rationale | Mitigation | Owner and review condition |
| --- | --- | --- | --- |
| A caller-supplied HTTP reader can block despite cancellation | Arbitrary `io.Reader.Read` has no interrupt contract; adding a hidden goroutine would leak work rather than guarantee interruption. | Use context-aware request bodies/readers and transport deadlines; core bounds reads and checks cancellation before/after. | Brian and application transport owner; review when reader support or cancellation semantics change, or blocking behavior is reported. |
| Collaborator error text can expose sensitive data | Validator/cache/resolver implementations are trusted application collaborators; this family cannot promise their errors are value-free. | Classify `errors.Is`/`errors.As` into fixed codes before observability; independently review collaborator redaction. | Brian and application observability owner; review on new collaborator/error propagation or disclosure report. |
| Trusted metadata may still be used without authentication/authorization | A conversion library cannot determine transport or application identity policy. | Authenticate the transport and authorize tenant/resource access before the explicit trusted extraction path; never use producer assertions alone. | Brian and application security owner; review on trust-policy/API changes or an authorization incident. |
| Returned event/envelope values can be exported by callers | Returning canonical application data is the library's purpose; owned copies are not redaction. | Deny raw values in default logs/traces/baggage/metrics; use finite allowlists, retention and access-control policy. | Brian and application data owner; review on new returned fields, observability integrations or disclosure reports. |

Maintainer compromise and malicious dependencies remain supply-chain risks:
Brian owns review of action/tool/dependency changes, vulnerability/license
results, key rotation and release incident response. Signatures identify the
signer; they do not independently prove source correctness.

## Verification and remaining qualification

Ordinary required CI, release rehearsals, actual public consumers and the native
family security diagnostic are distinct evidence. The diagnostic uses an
immutable development tooling source for owned analysis and full-history/current
secret checks across all declared modules; it is not a published Tools v2 claim
and cannot satisfy ordinary Required or release CI by itself. Scanner success
does not replace the source-level trust/privacy review above. Current results and
remaining adoption/release boundaries belong in the coordination ledger; this
model does not assert that all goal gates are already green.

Private reports must use the repository's private security reporting channel,
without production payloads or credentials unless a secure channel is agreed.
