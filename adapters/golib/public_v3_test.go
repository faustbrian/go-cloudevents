package golib_test

import (
	"testing"

	"github.com/faustbrian/go-cloudevents"
	golib "github.com/faustbrian/go-cloudevents/adapters/golib/v3"
	eventsourcing "github.com/faustbrian/go-event-sourcing/v2"
	registry "github.com/faustbrian/go-schema-registry/v2"
	"github.com/faustbrian/go-tenancy/v2"
)

// Public nominal identity is the adoption contract, not an implementation detail.
var _ func(eventsourcing.Message, golib.EventSourcingOptions) (cloudevents.Event, golib.EventSourcingState, golib.Report, error) = golib.EventSourcingToCloudEvent
var _ map[string]registry.Lookup = golib.RegistryJSONSchemaConfig{}.SchemaLookups
var _ tenancy.TenantID = golib.AuditMetadata{}.Tenant

func TestPublicTenantV2FacadeRoundTrip(t *testing.T) {
	data, err := cloudevents.NewJSONData([]byte(`{"ok":true}`))
	if err != nil {
		t.Fatal(err)
	}
	event, err := cloudevents.NewEvent(cloudevents.Attributes{ID: "event-1", Source: "/orders", Type: "order.created"}, data)
	if err != nil {
		t.Fatal(err)
	}
	tenant := tenancy.MustTenantID("tenant-a")
	converted, err := golib.AddTenant(event, tenant)
	if err != nil {
		t.Fatal(err)
	}
	var extracted tenancy.TenantID
	extracted, err = golib.ExtractTenant(converted, true)
	if err != nil || !extracted.Equal(tenant) {
		t.Fatalf("tenant round trip = %v, %v", extracted, err)
	}
	if _, present := event.Extension("tenantid"); present {
		t.Fatal("conversion modified the input event")
	}
}
