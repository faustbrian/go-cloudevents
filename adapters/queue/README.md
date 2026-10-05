# CloudEvents queue adapter

This module maps Golib queue jobs to CloudEvents while retaining retry,
settlement, execution, and operational state outside the event.

Install the published `adapters/queue/v2.0.0` release with
`go get github.com/faustbrian/go-cloudevents/adapters/queue/v2@v2.0.0`.
Use `ToCloudEvent` and retain the returned job; pass it to `FromCloudEvent` to
reconstruct the canonical queue value. The application owns queue I/O,
acknowledgements, retry policy, and lifecycle.

## V2 migration

Move queue-adapter imports to `/v2` and match tenant validation errors against
`github.com/faustbrian/go-tenancy/v2` sentinels. This error-identity boundary
uses an independent adapter major even though `job.Message` tenant metadata
remains string-valued. Queue core types, CloudEvents, retained execution state,
trust, ownership and loss reports are unchanged. Published adapter v1 keeps
its original Tenancy validation dependency and error identities.

Source stays in `adapters/queue/` on main; no version-specific directory or
root-module major is introduced. `api/baseline.txt` preserves v1;
`api/v2-baseline.txt` describes the new major. Released facade v1 and v2 keep
their original contracts; facade-v3 source and target integration now consume
this public v2 adapter.

See the [API reference](https://pkg.go.dev/github.com/faustbrian/go-cloudevents/adapters/queue/v2),
[parent documentation](https://github.com/faustbrian/go-cloudevents/blob/main/docs/README.md),
and [security policy](https://github.com/faustbrian/go-cloudevents/security/policy).
