package telemetry

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/faustbrian/go-cloudevents"
	telemetrypropagation "github.com/faustbrian/go-telemetry/propagation"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"
)

func TestInjectTraceContextPreservesOwnershipAndReportsLoss(t *testing.T) {
	t.Parallel()

	policy := propagationPolicy(t, true)
	event := telemetryEvent(t, nil)
	var nilContext context.Context
	if _, _, err := InjectTraceContext(nilContext, event, policy); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("nil context error = %v", err)
	}
	if _, _, err := InjectTraceContext(context.Background(), event, nil); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("nil policy error = %v", err)
	}
	if _, _, err := InjectTraceContext(telemetryTraceContext(t), cloudevents.Event{}, policy); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("zero event error = %v", err)
	}

	member, err := baggage.NewMember("tenant", "tenant-a")
	if err != nil {
		t.Fatal(err)
	}
	bag, err := baggage.New(member)
	if err != nil {
		t.Fatal(err)
	}
	ctx := baggage.ContextWithBaggage(telemetryTraceContext(t), bag)
	converted, report, err := InjectTraceContext(ctx, event, policy)
	if err != nil {
		t.Fatal(err)
	}
	if got := report.Losses; !reflect.DeepEqual(got, []Loss{{Field: "baggage", Reason: "no selected CloudEvents extension"}}) {
		t.Fatalf("losses = %#v", got)
	}
	for _, name := range []string{"traceparent", "tracestate"} {
		if _, ok := converted.Extension(name); !ok {
			t.Fatalf("%s extension is absent", name)
		}
		if _, ok := event.Extension(name); ok {
			t.Fatalf("input event gained %s", name)
		}
	}

	existing, err := cloudevents.NewTraceParentAttribute("00-aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa-bbbbbbbbbbbbbbbb-01")
	if err != nil {
		t.Fatal(err)
	}
	collision := telemetryEvent(t, map[string]cloudevents.Attribute{"traceparent": existing})
	if _, _, err := InjectTraceContext(ctx, collision, policy); !errors.Is(err, ErrMetadataCollision) {
		t.Fatalf("collision error = %v", err)
	}
}

func TestExtractTraceContextAppliesCallerTrustDecision(t *testing.T) {
	t.Parallel()

	var nilContext context.Context
	if got := ExtractTraceContext(nilContext, telemetryEvent(t, nil), nil, true); got == nil {
		t.Fatal("nil context was not replaced")
	}

	policy := propagationPolicy(t, false)
	event, _, err := InjectTraceContext(telemetryTraceContext(t), telemetryEvent(t, nil), policy)
	if err != nil {
		t.Fatal(err)
	}
	for _, trusted := range []bool{false, true} {
		extracted := ExtractTraceContext(context.Background(), event, policy, trusted)
		got := trace.SpanContextFromContext(extracted)
		if !got.IsValid() || got.TraceID() != telemetrySpanContext(t).TraceID() {
			t.Fatalf("trusted=%v span context = %v", trusted, got)
		}
	}

	withoutExtensions := ExtractTraceContext(context.Background(), telemetryEvent(t, nil), policy, false)
	if got := trace.SpanContextFromContext(withoutExtensions); got.IsValid() {
		t.Fatalf("absent trace extensions produced %v", got)
	}
	if got := trace.SpanContextFromContext(ExtractTraceContext(context.Background(), cloudevents.Event{}, policy, false)); got.IsValid() {
		t.Fatalf("zero event produced %v", got)
	}
}

func TestMetadataCarrierImplementsDeterministicTextMapCarrier(t *testing.T) {
	t.Parallel()

	carrier := metadataCarrier{}
	carrier.Set("tracestate", "vendor=value")
	carrier.Set("traceparent", "parent")
	if got := carrier.Get("traceparent"); got != "parent" {
		t.Fatalf("traceparent = %q", got)
	}
	if got := carrier.Keys(); !reflect.DeepEqual(got, []string{"traceparent", "tracestate"}) {
		t.Fatalf("keys = %#v", got)
	}
}

func propagationPolicy(t *testing.T, baggageEnabled bool) *telemetrypropagation.Policy {
	t.Helper()
	policy, err := telemetrypropagation.New(telemetrypropagation.Config{
		BaggageEnabled:     baggageEnabled,
		TrustedBaggageKeys: []string{"tenant"},
		MaxHeaderBytes:     1024,
		MaxBaggageItems:    1,
	})
	if err != nil {
		t.Fatal(err)
	}
	return policy
}

func telemetryTraceContext(t *testing.T) context.Context {
	t.Helper()
	return trace.ContextWithRemoteSpanContext(context.Background(), telemetrySpanContext(t))
}

func telemetrySpanContext(t *testing.T) trace.SpanContext {
	t.Helper()
	state, err := trace.ParseTraceState("vendor=value")
	if err != nil {
		t.Fatal(err)
	}
	return trace.NewSpanContext(trace.SpanContextConfig{
		TraceID:    trace.TraceID{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16},
		SpanID:     trace.SpanID{1, 2, 3, 4, 5, 6, 7, 8},
		TraceFlags: trace.FlagsSampled,
		TraceState: state,
		Remote:     true,
	})
}

func telemetryEvent(t *testing.T, extensions map[string]cloudevents.Attribute) cloudevents.Event {
	t.Helper()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "event-1", Source: "/source", Type: "example.created", Extensions: extensions,
	}, cloudevents.NewBinaryData([]byte("body")))
	if err != nil {
		t.Fatal(err)
	}
	return event
}
