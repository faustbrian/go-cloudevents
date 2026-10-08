// Package golib provides explicit, loss-reporting conversions between
// CloudEvents and Golib's canonical envelopes and metadata contracts.
//
// Conversions are synchronous, perform no implicit I/O, and preserve richer
// canonical state through explicit retained values. Inbound correlation,
// tenancy, tracing, and audit metadata is adopted only after a caller-owned
// trust decision. Queue and outbox adapters are Golib mappings, not official
// CloudEvents protocol bindings.
//
// Deprecated: use only the target-oriented adapters required by the
// application boundary. Version 4 adopts public JSON Schema v2, Outbox v2,
// and Schema Registry v3 alongside Tenancy v2, EventSourcing v2 and Workflow v2.
// Released v1, v2 and v3 remain
// available. All generations are excluded from the recommended adapter set
// because the facade owns the complete 26-dependency bridge.
package golib
