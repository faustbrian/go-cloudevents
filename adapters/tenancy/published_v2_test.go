package tenancy

import (
	"errors"
	"testing"

	"github.com/faustbrian/go-cloudevents"
	published "github.com/faustbrian/go-tenancy/v2"
)

func TestPublishedV2TenantIdentityRoundTrip(t *testing.T) {
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "event-1", Source: "/orders", Type: "example.created",
	}, cloudevents.NewBinaryData([]byte("body")))
	if err != nil {
		t.Fatal(err)
	}
	tenant := published.MustTenantID("tenant-a")
	converted, err := Add(event, tenant)
	if err != nil {
		t.Fatal(err)
	}
	extracted, err := Extract(converted, true)
	if err != nil || !extracted.Equal(tenant) {
		t.Fatalf("published identity round trip failed: %v", err)
	}
	if _, present := event.Extension("tenantid"); present {
		t.Fatal("input event mutated")
	}
	if missing, err := Extract(event, true); missing.Valid() || !errors.Is(err, published.ErrTenantMetadataMissing) {
		t.Fatal("missing metadata did not retain the published v2 sentinel")
	}
}
