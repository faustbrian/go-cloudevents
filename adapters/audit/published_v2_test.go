package audit_test

import (
	"testing"

	cloudaudit "github.com/faustbrian/go-cloudevents/adapters/audit/v2"
	published "github.com/faustbrian/go-tenancy/v2"
)

func TestPublishedV2TenantIdentityRoundTrip(t *testing.T) {
	want := published.MustTenantID("tenant-a")
	original := baseEvent(t)
	event, _, err := cloudaudit.AddMetadata(original, auditRecord(t, want.Value(), "", ""))
	if err != nil {
		t.Fatal(err)
	}
	metadata, err := cloudaudit.ExtractMetadata(event, true)
	if err != nil || !metadata.Tenant.Equal(want) {
		t.Fatalf("published tenant round trip failed: %v", err)
	}
	if _, present := original.Extension("tenantid"); present {
		t.Fatal("input event mutated")
	}
}
