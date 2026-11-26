# Changelog

## Unreleased

## 2.1.1 - 2026-11-26

### Security

- Require `github.com/klauspost/compress` v1.20.2 (from v1.19.2; GO-2026-6865),
  `golang.org/x/crypto` v0.59.0 (from v0.55.0; GO-2026-5932, GO-2026-6354 and
  GO-2026-6355) and `golang.org/x/net` v0.62.0 (from v0.58.0; GO-2026-6603,
  GO-2026-6610, GO-2026-6611, GO-2026-6612, GO-2026-6617, GO-2026-6860 and
  GO-2026-6866). The vulnerable versions were reachable only from tests.

## 2.0.1 - 2026-10-07

### Changed

- Adopt Correlation v1.1.2 and its indirect Identifier v2 dependency
  while retaining the audit adapter v2 correlation values, selected metadata,
  explicit trust, and loss-report contracts.

- Adopt Tenancy v2.0.1 in the adapter dependency graph while retaining
  the public v2 tenant identity.

### Documentation

- Record public adapter-v2 availability and explicit facade-v3 and target
  integration adoption; released older facade contracts remain available.

## 2.0.0 - 2026-10-05

- Adopt public Tenancy v2 identities and the public Tenancy adapter v2 through
  the independent `adapters/audit/v2` module. Migrate audit-adapter imports and
  `Metadata.Tenant` to `go-tenancy/v2` together; wrapped tenant errors now belong
  to that major. Mapping, explicit trust, ownership and loss reports are unchanged.
- Retain the released v1 API snapshot; published v1 callers keep their original
  tenant types and dependencies.

## 1.0.0 - 2026-09-09

- Publish the target-oriented audit metadata adapter with explicit trust,
  ownership, collision, and conversion-loss behavior.
