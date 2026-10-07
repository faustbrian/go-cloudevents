# Changelog

## Unreleased

## 2.0.1 - 2026-10-07

### Changed

- Adopt Queue v1.1.2 while retaining the public adapter v2 job
  conversion contract.

- Adopt Tenancy v2.0.1 in the adapter dependency graph while retaining
  the public v2 tenant identity.

### Documentation

- Record public adapter-v2 availability and explicit facade-v3 and target
  integration adoption; released older facade contracts remain available.

## 2.0.0 - 2026-10-05

- Adopt public Tenancy v2 validation through the independent `adapters/queue/v2`
  module. Migrate adapter imports and tenant-error matching to `go-tenancy/v2`
  together; queue jobs retain string tenant metadata and their existing core
  types. Round trips, caller-owned state, collisions and loss reports are unchanged.
- Retain the released v1 API snapshot and public v1 dependency line.

## 1.0.0 - 2026-09-09

- Publish the target-oriented queue adapter with retained execution state,
  nil-body preservation, trust, collision, and loss behavior.
