# Changelog

## Unreleased

### Changed

- Adopt Correlation v1.1.2 and its indirect Identifier v2 dependency
  while retaining identifier values, parsing policy, explicit trust,
  collision handling, and caller-owned event data.

- Align the indirect PostgreSQL client graph with pgx v5.11.0 while retaining
  the adapter's trust, validation, collision, and copy contracts.

## 1.0.0 - 2026-09-09

- Publish the target-oriented correlation adapter with explicit trust,
  validation, collision, and copy semantics.
