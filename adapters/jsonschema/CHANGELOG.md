# Changelog

## Unreleased

- Prepare the `/v2` adapter with public `go-json-schema/v2` v2.0.0 compiled
  schema types. Upgrade adapter and schema imports together; the CloudEvents
  root remains v1. Preserve exact URI admission, parser and cancellation
  errors, Unicode code-point validation and schema-violation classification.

### Security

- Require `github.com/klauspost/compress` v1.20.2 (from v1.19.2; GO-2026-6865),
  `golang.org/x/crypto` v0.59.0 (from v0.57.0; GO-2026-5932) and
  `golang.org/x/net` v0.62.0 (from v0.59.0; GO-2026-6603, GO-2026-6610,
  GO-2026-6611, GO-2026-6612, GO-2026-6617, GO-2026-6860 and GO-2026-6866). The
  vulnerable `github.com/klauspost/compress` and `golang.org/x/crypto` versions
  were reachable only from tests. Production packages import
  `golang.org/x/net`, but govulncheck found no imported vulnerable package.

## 1.0.0 - 2026-09-09

- Publish the target-oriented direct JSON Schema validator with explicit URI,
  content-type, context, and validation behavior.
