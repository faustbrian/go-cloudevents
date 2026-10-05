# Changelog

## Unreleased

### Documentation

- Record public adapter-v2 availability and explicit facade-v3 and target
  integration adoption; released older facade contracts remain available.

- Adopt public Tenancy v2 identities and the public Tenancy adapter v2 through
  the independent `adapters/audit/v2` module. Migrate audit-adapter imports and
  `Metadata.Tenant` to `go-tenancy/v2` together; wrapped tenant errors now belong
  to that major. Mapping, explicit trust, ownership and loss reports are unchanged.
- Retain the released v1 API snapshot; published v1 callers keep their original
  tenant types and dependencies.

## 1.0.0 - 2026-09-09

- Publish the target-oriented audit metadata adapter with explicit trust,
  ownership, collision, and conversion-loss behavior.
