# Changelog

## Unreleased

## 1.1.1 - 2026-11-26

### Security

- Require klauspost/compress v1.20.2 to address GO-2026-6865.

- Require `golang.org/x/crypto` v0.59.0 (from v0.55.0; GO-2026-5932,
  GO-2026-6354 and GO-2026-6355), `golang.org/x/net` v0.62.0 (from v0.58.0;
  GO-2026-6603, GO-2026-6610, GO-2026-6611, GO-2026-6612, GO-2026-6617,
  GO-2026-6860 and GO-2026-6866) and OpenTelemetry v1.47.0 (from v1.46.0;
  GO-2026-6974 and GO-2026-6975). The vulnerable `golang.org/x/net` and
  OpenTelemetry versions were reachable only from tests. Production packages
  import `golang.org/x/crypto`, but govulncheck found no imported vulnerable
  package.

## 1.0.2 - 2026-10-07

### Changed

- Select Kafka root module v1.1.0 while preserving record ownership,
  transport metadata, and CloudEvents binding behavior.

- Align indirect OpenTelemetry metric and trace modules with v1.46.0 and
  reconcile their checksums without changing record mapping or ownership.

## 1.0.0 - 2026-09-09

- Publish the target-oriented Kafka record adapter with bounded decoding,
  collision detection, and explicit transport ownership.
