# Changelog

## Unreleased

### Security

- Require `github.com/klauspost/compress` v1.20.2 (from v1.19.2; GO-2026-6865),
  `golang.org/x/crypto` v0.59.0 (from v0.55.0; GO-2026-5932, GO-2026-6354 and
  GO-2026-6355) and `golang.org/x/net` v0.62.0 (from v0.58.0; GO-2026-6603,
  GO-2026-6610, GO-2026-6611, GO-2026-6612, GO-2026-6617, GO-2026-6860 and
  GO-2026-6866). The vulnerable versions were reachable only from tests.

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
