# Changelog

## Unreleased

## 1.0.2 - 2026-10-07

### Changed

- Select Kafka root module v1.1.0 while preserving record ownership,
  transport metadata, and CloudEvents binding behavior.

- Align indirect OpenTelemetry metric and trace modules with v1.46.0 and
  reconcile their checksums without changing record mapping or ownership.

## 1.0.0 - 2026-09-09

- Publish the target-oriented Kafka record adapter with bounded decoding,
  collision detection, and explicit transport ownership.
