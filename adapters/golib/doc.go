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
// application boundary. Version 3 adopts public Tenancy v2, EventSourcing v2,
// and Schema Registry v2 types alongside Workflow v2. Released v1 and v2 remain
// available. All generations are excluded from the recommended adapter set
// because the facade owns the complete 26-dependency bridge.
package golib
