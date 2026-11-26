# Changelog

## Unreleased

### Security

- Require `github.com/klauspost/compress` v1.20.2 (from v1.19.2; GO-2026-6865),
  `golang.org/x/crypto` v0.59.0 (from v0.55.0; GO-2026-5932, GO-2026-6354 and
  GO-2026-6355) and `golang.org/x/net` v0.62.0 (from v0.58.0; GO-2026-6603,
  GO-2026-6610, GO-2026-6611, GO-2026-6612, GO-2026-6617, GO-2026-6860 and
  GO-2026-6866). The vulnerable versions were reachable only from tests.

## 2.0.1 - 2026-10-07

### Documentation

- Record public adapter-v2 availability and explicit facade-v3 and target
  integration adoption; released older facade contracts remain available.

### Changed

- Adopt Tenancy v2.0.1 in the adapter dependency graph while retaining
  the public v2 tenant identity.

## 2.0.0 - 2026-10-05

### Changed

- Use public Tenancy v2 validation in the pending Event Sourcing adapter v2;
  wrapped tenant errors now match `go-tenancy/v2` sentinels. Retained tenant
  state stays string-valued, with unchanged conversions, trust and ownership.
- Adopt public Event Sourcing core v2.0.0 through the independent
  `github.com/faustbrian/go-cloudevents/adapters/event-sourcing/v2` module
  identity. Move imports and core message/state types together; conversion
  behavior, retained state, ownership, and loss reporting remain unchanged.
- Require the existing public CloudEvents root v1.1.1 dependency release,
  retaining its unsuffixed identity and the Go 1.27.0 minimum.

### Documentation

- Preserve the released v1.0.1 API snapshot and distinguish pending v2
  publication from later composition-consumer adoption. The deprecated Golib
  compatibility facade remains on its released v1 contracts.

## 1.0.0 - 2026-09-09

- Publish the target-oriented event-sourcing adapter with explicit retained
  store state, trust, collision, and conversion-loss behavior.
