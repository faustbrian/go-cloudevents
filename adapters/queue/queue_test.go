package queue

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-queue/job"
	"github.com/faustbrian/go-tenancy"
)

func TestQueueRoundTripPreservesOwnedStateWithoutAliasing(t *testing.T) {
	t.Parallel()

	enqueuedAt := time.Date(2026, 8, 10, 1, 2, 3, 0, time.UTC)
	occurredAt := enqueuedAt.Add(time.Hour)
	message := job.Message{
		Timeout: time.Minute, Body: []byte("payload"), RetryCount: 2, RetryDelay: time.Second,
		Metadata: &job.Metadata{
			OriginalID: "job-1", JobType: "order.notify", ContentType: "application/octet-stream",
			EnqueuedAt: &enqueuedAt, TenantID: "tenant-a",
			Tags: map[string]string{"priority": "high"},
			Correlation: map[string]string{
				"correlationid": "correlation-1", "requestid": "request-1", "causationid": "cause-1",
			},
			TraceContext: map[string]string{
				"traceparent": "00-4bf92f3577b34da6a3ce929d0e0e4736-00f067aa0ba902b7-01",
				"tracestate":  "vendor=value",
			},
		},
	}
	event, retained, report, err := ToCloudEvent(message, Options{Source: "/queue", OccurredAt: &occurredAt})
	if err != nil || len(report.Losses) != 0 {
		t.Fatalf("ToCloudEvent() = %#v, %v", report, err)
	}
	message.Body[0] = 'X'
	message.Metadata.Tags["priority"] = "low"
	message.Metadata.Correlation["correlationid"] = "changed"
	message.Metadata.TraceContext["tracestate"] = "changed"
	enqueuedAt = time.Time{}
	occurredAt = time.Time{}
	if string(retained.Body) != "payload" || retained.Metadata.Tags["priority"] != "high" ||
		retained.Metadata.Correlation["correlationid"] != "correlation-1" ||
		retained.Metadata.TraceContext["tracestate"] != "vendor=value" || retained.Metadata.EnqueuedAt.IsZero() {
		t.Fatalf("retained state aliases input: %#v", retained)
	}

	roundTrip, reverseReport, err := FromCloudEvent(event, retained)
	if err != nil || !bytes.Equal(roundTrip.Body, []byte("payload")) || roundTrip.RetryCount != 2 ||
		roundTrip.Metadata.OriginalID != "job-1" || len(reverseReport.Losses) < 1 {
		t.Fatalf("FromCloudEvent() = %#v, %#v, %v", roundTrip, reverseReport, err)
	}
	roundTrip.Body[0] = 'Y'
	roundTrip.Metadata.Tags["priority"] = "urgent"
	if string(retained.Body) != "payload" || retained.Metadata.Tags["priority"] != "high" {
		t.Fatal("round-trip state aliases retained state")
	}
}

func TestQueueConversionRejectsInvalidInputsAndCollisions(t *testing.T) {
	t.Parallel()

	valid := job.Message{Timeout: time.Second, Body: []byte("body")}
	invalidTenant := valid
	invalidTenant.Metadata = &job.Metadata{TenantID: "bad?tenant", OriginalID: "id", JobType: "type"}
	if _, _, _, err := ToCloudEvent(invalidTenant, Options{Source: "/source"}); !errors.Is(err, tenancy.ErrInvalidTenantID) {
		t.Fatalf("invalid tenant error = %v", err)
	}
	for name, input := range map[string]struct {
		message job.Message
		options Options
	}{
		"invalid message":  {message: job.Message{}, options: Options{Source: "/source", StableID: "id", Type: "type"}},
		"missing source":   {message: valid, options: Options{StableID: "id", Type: "type"}},
		"missing identity": {message: valid, options: Options{Source: "/source"}},
		"missing ID":       {message: valid, options: Options{Source: "/source", Type: "type"}},
		"missing type":     {message: valid, options: Options{Source: "/source", StableID: "id"}},
	} {
		if _, _, _, err := ToCloudEvent(input.message, input.options); !errors.Is(err, ErrInvalidInput) {
			t.Fatalf("%s error = %v", name, err)
		}
	}
	invalidJSON := valid
	invalidJSON.Metadata = &job.Metadata{OriginalID: "id", JobType: "type", ContentType: "application/json"}
	invalidJSON.Body = []byte("{")
	if _, _, _, err := ToCloudEvent(invalidJSON, Options{Source: "/source"}); err == nil {
		t.Fatal("invalid JSON error = nil")
	}
	invalidExtension := valid
	invalidExtension.Metadata = &job.Metadata{OriginalID: "id", JobType: "type", Correlation: map[string]string{"correlationid": "bad\nvalue"}}
	if _, _, _, err := ToCloudEvent(invalidExtension, Options{Source: "/source"}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("invalid extension error = %v", err)
	}
	if _, _, _, err := ToCloudEvent(valid, Options{Source: "bad\nsource", StableID: "id", Type: "type"}); err == nil {
		t.Fatal("invalid event source error = nil")
	}

	event := queueEvent(t, "id", "type", "", nil, cloudevents.NewBinaryData([]byte("body")))
	if _, _, err := FromCloudEvent(cloudevents.Event{}, valid); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("invalid event error = %v", err)
	}
	if _, _, err := FromCloudEvent(queueEvent(t, "id", "type", "", nil, cloudevents.Data{}), valid); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("absent data error = %v", err)
	}
	if state, _, err := FromCloudEvent(event, valid); err != nil || state.Metadata == nil {
		t.Fatalf("nil metadata restoration = %#v, %v", state, err)
	}
	for _, state := range []job.Message{
		{Timeout: time.Second, Metadata: &job.Metadata{OriginalID: "other"}},
		{Timeout: time.Second, Metadata: &job.Metadata{JobType: "other"}},
		{Timeout: time.Second, Metadata: &job.Metadata{OriginalID: "id", JobType: "type", ContentType: "text/plain"}},
	} {
		if _, _, err := FromCloudEvent(event, state); !errors.Is(err, ErrMetadataCollision) {
			t.Fatalf("metadata collision error = %v", err)
		}
	}
}

func TestQueueExtensionValidationAndLossReporting(t *testing.T) {
	t.Parallel()

	metadata := &job.Metadata{
		OriginalID: "id", JobType: "type", TenantID: "tenant-a",
		Correlation: map[string]string{"correlationid": "correlation-1"},
	}
	extensions, err := queueExtensions(metadata)
	if err != nil || len(extensions) != 2 {
		t.Fatalf("queueExtensions() = %#v, %v", extensions, err)
	}
	if empty, err := queueExtensions(nil); err != nil || len(empty) != 0 {
		t.Fatalf("queueExtensions(nil) = %#v, %v", empty, err)
	}
	if _, err := queueExtensions(&job.Metadata{TenantID: "bad?tenant"}); !errors.Is(err, tenancy.ErrInvalidTenantID) {
		t.Fatalf("invalid tenant extension error = %v", err)
	}

	base := job.Message{Timeout: time.Second, Body: []byte("body"), Metadata: metadata}
	event, retained, _, err := ToCloudEvent(base, Options{Source: "/source"})
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"correlationid", "tenantid"} {
		changed := event.Extensions()
		delete(changed, name)
		if _, _, err := FromCloudEvent(queueEvent(t, event.ID(), event.Type(), "", changed, event.Data()), retained); !errors.Is(err, ErrMetadataCollision) {
			t.Fatalf("missing %s collision error = %v", name, err)
		}
	}
	other, _ := cloudevents.NewStringAttribute("other")
	changed := event.Extensions()
	changed["correlationid"] = other
	if _, _, err := FromCloudEvent(queueEvent(t, event.ID(), event.Type(), "", changed, event.Data()), retained); !errors.Is(err, ErrMetadataCollision) {
		t.Fatalf("changed extension collision error = %v", err)
	}
	changed = event.Extensions()
	changed["correlationid"] = cloudevents.NewBooleanAttribute(true)
	if _, _, err := FromCloudEvent(queueEvent(t, event.ID(), event.Type(), "", changed, event.Data()), retained); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("wrong extension kind error = %v", err)
	}
	badTenant, _ := cloudevents.NewStringAttribute("bad?tenant")
	changed = event.Extensions()
	changed["tenantid"] = badTenant
	if _, _, err := FromCloudEvent(queueEvent(t, event.ID(), event.Type(), "", changed, event.Data()), retained); !errors.Is(err, tenancy.ErrInvalidTenantID) {
		t.Fatalf("wire tenant error = %v", err)
	}
	retained.Metadata.TenantID = "bad?tenant"
	if _, _, err := FromCloudEvent(event, retained); !errors.Is(err, tenancy.ErrInvalidTenantID) {
		t.Fatalf("retained tenant error = %v", err)
	}

	extra, _ := cloudevents.NewStringAttribute("value")
	jsonData, _ := cloudevents.NewJSONData([]byte(`{"ok":true}`))
	lossEvent, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "id", Source: "/source", Type: "type", Subject: "subject",
		DataContentType: "application/json", DataSchema: "https://example.test/schema",
		Extensions: map[string]cloudevents.Attribute{"extra": extra}, Time: timePointer(time.Now()),
	}, jsonData)
	if err != nil {
		t.Fatal(err)
	}
	state := job.Message{Timeout: time.Second, Metadata: &job.Metadata{OriginalID: "id", JobType: "type", ContentType: "application/json"}}
	if _, report, err := FromCloudEvent(lossEvent, state); err != nil || len(report.Losses) < 5 {
		t.Fatalf("loss report = %#v, %v", report, err)
	}
}

func TestQueuePreservesNilAndEmptyBodies(t *testing.T) {
	t.Parallel()

	for name, body := range map[string][]byte{"nil": nil, "empty": {}} {
		body := body
		t.Run(name, func(t *testing.T) {
			message := job.Message{Timeout: time.Second, Body: body, Metadata: &job.Metadata{OriginalID: "id", JobType: "type"}}
			event, retained, _, err := ToCloudEvent(message, Options{Source: "/queue"})
			if err != nil {
				t.Fatal(err)
			}
			roundTrip, _, err := FromCloudEvent(event, retained)
			if err != nil || (body == nil) != (roundTrip.Body == nil) || len(roundTrip.Body) != 0 {
				t.Fatalf("round-trip body = %#v, %v", roundTrip.Body, err)
			}
		})
	}
}

func TestQueuePrivateHelpersCoverZeroAndCopiedState(t *testing.T) {
	t.Parallel()

	message := job.Message{Body: []byte("body")}
	if cloned := cloneJob(message); string(cloned.Body) != "body" || cloned.Metadata != nil {
		t.Fatalf("cloneJob() = %#v", cloned)
	}
	if value, err := validatedTenant(""); err != nil || value != "" {
		t.Fatalf("validatedTenant(empty) = %q, %v", value, err)
	}
	if value, err := validatedTenant("tenant-a"); err != nil || value != "tenant-a" {
		t.Fatalf("validatedTenant(valid) = %q, %v", value, err)
	}
}

func queueEvent(t *testing.T, id, eventType, contentType string, extensions map[string]cloudevents.Attribute, data cloudevents.Data) cloudevents.Event {
	t.Helper()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: id, Source: "/source", Type: eventType, DataContentType: contentType, Extensions: extensions,
	}, data)
	if err != nil {
		t.Fatal(err)
	}
	return event
}

func timePointer(value time.Time) *time.Time { return &value }
