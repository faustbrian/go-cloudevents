// Package cloudevents implements the transport-independent CloudEvents 1.0
// information model together with explicitly selected event formats, protocol
// bindings, and extensions. The pinned base edition is CloudEvents 1.0.2;
// the wire specversion remains "1.0". Selected stable-line clarifications and
// defensive interpretations are recorded in docs/specification-decisions.md,
// with immutable source identities in specification/manifest.json.
// NewEvent is the value entry point; DecodeJSON/HTTP/Kafka admit wire inputs.
// Constructors/encoders do not apply Limits, authenticate sources or provide
// canonical signing bytes. Decoders apply selected byte/count/depth checks,
// not a whole-process heap/CPU budget. Copies are ordinary plaintext storage,
// not redaction or erasure. Immutable Event/Attribute/Data values may be reused;
// caller-owned maps/readers/callbacks must remain valid during each operation.
//
// The package is an interoperability envelope, not an event bus, event store,
// queue, outbox, workflow engine, schema registry, or canonical domain-event
// model. Constructors and decoders own retained mutable input. The core does no
// network I/O, schema resolution, global registration, telemetry, or background
// work.
package cloudevents
