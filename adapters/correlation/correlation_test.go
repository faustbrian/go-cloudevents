package correlation_test

import (
	"errors"
	"testing"

	"github.com/faustbrian/go-cloudevents"
	cloudcorrelation "github.com/faustbrian/go-cloudevents/adapters/correlation"
	"github.com/faustbrian/go-correlation"
)

func TestCorrelationRoundTripPreservesAllIdentifiers(t *testing.T) {
	t.Parallel()

	values := correlation.Values{
		CorrelationID: correlation.MustCorrelationID("correlation-1", correlation.Policy{}),
		RequestID:     correlation.MustRequestID("request-1", correlation.Policy{}),
		CausationID:   correlation.MustCausationID("cause-1", correlation.Policy{}),
	}
	event, err := cloudcorrelation.Add(baseEvent(t), values)
	if err != nil {
		t.Fatal(err)
	}
	extracted, err := cloudcorrelation.Extract(event, true, correlation.Policy{})
	if err != nil || extracted != values {
		t.Fatalf("Extract() = %#v, %v; want %#v, nil", extracted, err, values)
	}
	if _, err := cloudcorrelation.Add(event, values); err != nil {
		t.Fatalf("idempotent Add() error = %v", err)
	}
}

func TestCorrelationAbsentMetadataReturnsZeroWithoutTrust(t *testing.T) {
	t.Parallel()

	values, err := cloudcorrelation.Extract(baseEvent(t), false, correlation.Policy{})
	if err != nil || values != (correlation.Values{}) {
		t.Fatalf("Extract() = %#v, %v; want zero values, nil", values, err)
	}
}

func TestCorrelationRejectsUntrustedCollidingAndMalformedMetadata(t *testing.T) {
	t.Parallel()

	valid := correlation.Values{CorrelationID: correlation.MustCorrelationID("correlation-1", correlation.Policy{})}
	event, err := cloudcorrelation.Add(baseEvent(t), valid)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := cloudcorrelation.Extract(event, false, correlation.Policy{}); !errors.Is(err, cloudcorrelation.ErrUntrustedMetadata) {
		t.Fatalf("untrusted Extract() error = %v", err)
	}
	if _, err := cloudcorrelation.Add(event, correlation.Values{CorrelationID: "different"}); !errors.Is(err, cloudcorrelation.ErrMetadataCollision) {
		t.Fatalf("colliding Add() error = %v", err)
	}
	if _, err := cloudcorrelation.Add(cloudevents.Event{}, correlation.Values{}); !errors.Is(err, cloudcorrelation.ErrInvalidInput) {
		t.Fatalf("zero-event Add() error = %v", err)
	}

	for _, name := range []string{"correlationid", "requestid", "causationid"} {
		name := name
		t.Run(name+" wrong attribute kind", func(t *testing.T) {
			t.Parallel()
			event := eventWithExtensions(t, map[string]cloudevents.Attribute{name: cloudevents.NewBooleanAttribute(true)})
			if _, err := cloudcorrelation.Extract(event, true, correlation.Policy{}); !errors.Is(err, cloudcorrelation.ErrInvalidInput) {
				t.Fatalf("Extract() error = %v", err)
			}
		})
		t.Run(name+" invalid identifier", func(t *testing.T) {
			t.Parallel()
			attribute, err := cloudevents.NewStringAttribute("bad.value")
			if err != nil {
				t.Fatal(err)
			}
			event := eventWithExtensions(t, map[string]cloudevents.Attribute{name: attribute})
			if _, err := cloudcorrelation.Extract(event, true, correlation.Policy{}); !errors.Is(err, cloudcorrelation.ErrInvalidInput) {
				t.Fatalf("Extract() error = %v", err)
			}
		})
	}
}

func baseEvent(t *testing.T) cloudevents.Event {
	t.Helper()
	return eventWithExtensions(t, nil)
}

func eventWithExtensions(t *testing.T, extensions map[string]cloudevents.Attribute) cloudevents.Event {
	t.Helper()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "event-1", Source: "/source", Type: "example.created", Extensions: extensions,
	}, cloudevents.NewBinaryData([]byte("body")))
	if err != nil {
		t.Fatal(err)
	}
	return event
}
