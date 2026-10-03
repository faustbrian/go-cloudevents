# CloudEvents workflow adapter

This module maps Golib durable workflow history to CloudEvents while retaining
workflow-owned sequence, definition, retry, compensation, and scheduling
state outside the event.

After `adapters/workflow/v2.0.0` is publicly released, install with
`go get github.com/faustbrian/go-cloudevents/adapters/workflow/v2@v2.0.0`.
Use `ToCloudEvent` and retain its `State`; pass that state to `FromCloudEvent`
when reconstructing the canonical history event. Runtime lifecycle is external.

## Migration from v1

Version 2 adopts public `github.com/faustbrian/go-workflow/v2@v2.0.0`.
Update the adapter import and workflow imports together: `ToCloudEvent` accepts
a workflow-v2 `HistoryEvent`, `FromCloudEvent` returns one, and
`State.Definition` uses a workflow-v2 `DefinitionReference`. The v1 and v2
nominal Go types are not interchangeable.

CloudEvents remains at its unsuffixed root v1.1.0. Stable identity, sequence,
definition and retained state, occurrence time, payload copies and nil-versus-empty
payload preservation, error classifications and loss reports are unchanged.
The adapter adds no persistence, implicit I/O or workflow orchestration.

Source stays in `adapters/workflow/` on main; its independent release tag is
`adapters/workflow/v2.0.0`. No root-module release or version-specific source
directory is required. `api/v1.0.1.txt` preserves the exact released v1 API from
commit `5559c521abeb0ac979a06a9405185b8486033c88`; `api/baseline.txt` describes
the current major, not compatibility between different module identities.

Published adapter v1 remains available. The deprecated `adapters/golib` facade
and non-releasable `integration/target-adapters` composition intentionally remain
on released adapter-v1 and workflow-v1 until this producer is publicly released.
They require a separate deliberate import and dependency migration afterward;
the facade's public workflow types require its own major version.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/workflow/v2),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
