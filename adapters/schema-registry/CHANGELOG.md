# Changelog

## Unreleased

## 3.1.1 - 2026-11-26

### Security

- Require `github.com/klauspost/compress` v1.20.2 (from v1.19.2; GO-2026-6865),
  `golang.org/x/crypto` v0.59.0 (from v0.57.0; GO-2026-5932) and
  `golang.org/x/net` v0.62.0 (from v0.59.0; GO-2026-6603, GO-2026-6610,
  GO-2026-6611, GO-2026-6612, GO-2026-6617, GO-2026-6860 and GO-2026-6866). The
  vulnerable `github.com/klauspost/compress` and `golang.org/x/crypto` versions
  were reachable only from tests. Production packages import
  `golang.org/x/net`, but govulncheck found no imported vulnerable package.

## 3.1.0 - 2026-10-22

### Documentation

- Clarify the shared resolution/payload-validation timeout context and borrowed
  cache/adapter lifecycle; deadlines require cooperative collaborators.

## 3.0.0 - 2026-10-08

- Prepare the `/v3` adapter with public `go-schema-registry/v3` v3.0.0 cache,
  lookup, availability and JSON Schema adapter types. Upgrade both adapter
  and registry imports together. Keep exact URI admission, bounded resolution,
  private payload diagnostics and CloudEvents error classifications.

## 2.0.0 - 2026-10-02

- Adopt schema-registry v2 through the `/adapters/schema-registry/v2` module
  path. Configuration now accepts registry v2 cache, lookup, availability, and
  JSON Schema adapter types; update both adapter and registry imports together.
- Preserve explicit bounded validation and CloudEvents error classifications
  while using registry v2's private JSON payload diagnostics. The CloudEvents
  root and deprecated v1 facade retain their existing module paths.

## 1.0.0 - 2026-09-09

- Publish the target-oriented registry validator with static schema selection,
  bounded resolution, explicit availability, and context propagation.
