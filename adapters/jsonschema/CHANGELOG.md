# Changelog

## Unreleased

- Prepare the `/v2` adapter with public `go-json-schema/v2` v2.0.0 compiled
  schema types. Upgrade adapter and schema imports together; the CloudEvents
  root remains v1. Preserve exact URI admission, parser and cancellation
  errors, Unicode code-point validation and schema-violation classification.

## 1.0.0 - 2026-09-09

- Publish the target-oriented direct JSON Schema validator with explicit URI,
  content-type, context, and validation behavior.
