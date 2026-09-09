// Package telemetry maps caller-owned Golib propagation policy through
// CloudEvents trace extensions without owning telemetry initialization.
package telemetry

import (
	"context"
	"fmt"
	"sort"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	telemetrypropagation "github.com/faustbrian/go-telemetry/propagation"
)

var (
	// ErrInvalidInput identifies a nil caller-owned context or propagation
	// policy. Callers can match it with errors.Is.
	ErrInvalidInput = cloudevents.ErrInvalidAdapterInput
	// ErrMetadataCollision identifies an event whose existing trace extension
	// conflicts with the value produced by the caller-owned policy. The input
	// event remains unchanged.
	ErrMetadataCollision = cloudevents.ErrMetadataCollision
)

// Loss describes one propagation field that was deliberately not represented
// in CloudEvents metadata. Its zero value has no field or reason.
type Loss = cloudevents.AdapterLoss

// Report describes lossy aspects of a completed conversion. Its zero value
// means no loss was observed; callers own the returned value and its slice.
type Report = cloudevents.AdapterReport

// InjectTraceContext returns a new event containing traceparent and tracestate
// produced by the caller-owned policy. A nil context or policy returns
// ErrInvalidInput, as does the zero Event. Existing equal extensions are
// preserved, while conflicting extensions return ErrMetadataCollision without
// changing event. Baggage is never flattened into event extensions and is
// instead reported as a Loss.
func InjectTraceContext(ctx context.Context, event cloudevents.Event, policy *telemetrypropagation.Policy) (cloudevents.Event, Report, error) {
	if ctx == nil || policy == nil {
		return cloudevents.Event{}, Report{}, fmt.Errorf("%w: telemetry policy", ErrInvalidInput)
	}
	carrier := metadataCarrier{}
	policy.Inject(ctx, carrier)
	additions := map[string]string{}
	for _, name := range []string{"traceparent", "tracestate"} {
		if value := carrier.Get(name); value != "" {
			additions[name] = value
		}
	}
	converted, err := adapter.AddStringExtensions(event, additions)
	if err != nil {
		return cloudevents.Event{}, Report{}, err
	}
	report := Report{}
	if carrier.Get("baggage") != "" {
		report.Losses = []Loss{{Field: "baggage", Reason: "no selected CloudEvents extension"}}
	}
	return converted, report, nil
}

// ExtractTraceContext delegates inbound trust handling to the caller-owned
// policy. trusted selects trusted extraction; false applies the policy's
// untrusted filtering. A nil context is replaced with context.Background, and
// a nil policy returns that context unchanged. The event and policy remain
// caller-owned and are not mutated. A zero Event contributes no trace metadata.
func ExtractTraceContext(ctx context.Context, event cloudevents.Event, policy *telemetrypropagation.Policy, trusted bool) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	if policy == nil {
		return ctx
	}
	carrier := metadataCarrier{}
	for _, name := range []string{"traceparent", "tracestate"} {
		value, present, _ := adapter.StringExtension(event, name)
		if present {
			carrier.Set(name, value)
		}
	}
	if trusted {
		return policy.ExtractTrusted(ctx, carrier)
	}
	return policy.Extract(ctx, carrier)
}

type metadataCarrier map[string]string

func (carrier metadataCarrier) Get(key string) string { return carrier[key] }
func (carrier metadataCarrier) Set(key, value string) { carrier[key] = value }
func (carrier metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(carrier))
	for key := range carrier {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
