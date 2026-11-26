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

- Adopt Tenancy v2.0.1 while retaining the adapter and core public v2
  identities.

### Documentation

- Clarify installation of the already published adapter-v2 release.

## 2.0.0 - 2026-10-05

- Adopt published Tenancy v2 identities and errors under the adapter's v2
  module path. Update adapter and Tenancy imports together; wire metadata,
  explicit trust, collision handling and input ownership are unchanged.

## 1.0.0 - 2026-09-09

- Publish the target-oriented tenancy adapter with validation, explicit trust,
  and collision behavior.
