# Changelog

## Unreleased

- Prepare the `/v2` adapter with public transactional-outbox `/v2` v2.0.0
  envelope types. Upgrade adapter and outbox imports together; the CloudEvents
  root remains v1. Preserve retained-state ownership, nil-versus-empty payloads,
  identifier collision checks and explicit field-loss reporting.

## 1.0.0 - 2026-09-09

- Publish the target-oriented outbox adapter with retained relay state,
  nil-payload preservation, and explicit conversion losses.
