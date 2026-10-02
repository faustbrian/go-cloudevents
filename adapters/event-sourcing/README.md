# CloudEvents event-sourcing adapter

This module maps Golib event-sourcing messages to CloudEvents while keeping
stream versions, positions, persistence timestamps, and metadata in explicit
caller-owned state.

After `adapters/event-sourcing/v2.0.0` is publicly released, install with
`go get github.com/faustbrian/go-cloudevents/adapters/event-sourcing/v2@v2.0.0`.
Use `ToCloudEvent` and retain its `State`; pass that state to `FromCloudEvent`
when reconstructing the canonical message. No storage or transport is owned.

## V2 migration

This independent module adopts public
`github.com/faustbrian/go-event-sourcing/v2@v2.0.0` and the existing unsuffixed
CloudEvents root at v1.1.1, rather than v1.1.0. Tenancy remains v1.1.0.
Update adapter imports and core types together: `ToCloudEvent` accepts a
core-v2 `Message`, `FromCloudEvent` returns one, and `State` uses core-v2
`StreamID`, `GlobalPosition`, and `SchemaVersion`. Core-v1 and core-v2 named
types are not interchangeable. CloudEvents `Event`, loss/report aliases,
conversion behavior, and explicit ownership remain unchanged.

Source stays in `adapters/event-sourcing/` on main, with tag
`adapters/event-sourcing/v2.0.0` and Go 1.27.0. This is not a root-module major
or a new version-specific source directory. No PostgreSQL dependency is added.

Published adapter v1 remains available. The deprecated `adapters/golib` facade
retains its core-v1 and adapter-v1 contracts throughout its supported v1
interval; it does not expose this v2 adapter. The non-releasable
`integration/target-adapters` composition remains on v1 until public adapter-v2
publication, then requires an explicit import/dependency update.

`api/v1.0.1.txt` preserves the exact released API from commit
`5559c521abeb0ac979a06a9405185b8486033c88`; `api/baseline.txt` describes the
current major, not compatibility between distinct Go module identities.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/event-sourcing/v2),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
