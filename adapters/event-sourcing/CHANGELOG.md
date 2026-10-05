# Changelog

## Unreleased

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
