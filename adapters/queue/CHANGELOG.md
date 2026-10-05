# Changelog

## Unreleased

- Adopt public Tenancy v2 validation through the independent `adapters/queue/v2`
  module. Migrate adapter imports and tenant-error matching to `go-tenancy/v2`
  together; queue jobs retain string tenant metadata and their existing core
  types. Round trips, caller-owned state, collisions and loss reports are unchanged.
- Retain the released v1 API snapshot and public v1 dependency line.

## 1.0.0 - 2026-09-09

- Publish the target-oriented queue adapter with retained execution state,
  nil-body preservation, trust, collision, and loss behavior.
