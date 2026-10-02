# Changelog

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
