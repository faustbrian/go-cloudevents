# Changelog

## Unreleased

### Documentation

- Clarify that decoded State.Timestamp retains the transport message timestamp,
  while Encode supplies event occurrence time; neither is a broker-clock proof.

### Security

- Require `github.com/klauspost/compress` v1.20.2 (from v1.19.2; GO-2026-6865),
  `golang.org/x/crypto` v0.59.0 (from v0.55.0; GO-2026-5932, GO-2026-6354 and
  GO-2026-6355) and `golang.org/x/net` v0.62.0 (from v0.58.0; GO-2026-6603,
  GO-2026-6610, GO-2026-6611, GO-2026-6612, GO-2026-6617, GO-2026-6860 and
  GO-2026-6866). The vulnerable versions were reachable only from tests.

## 1.0.2 - 2026-10-07

### Changed

- Adopt RabbitMQ Streams v1.1.1 in the adapter dependency graph.

## 1.0.0 - 2026-09-09

- Publish the target-oriented RabbitMQ Streams adapter with structured JSON,
  bounded metadata copying, and explicit broker-state ownership.
