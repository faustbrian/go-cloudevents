package tenancy

import (
	"errors"
	"testing"

	"github.com/faustbrian/go-cloudevents"
	golibtenancy "github.com/faustbrian/go-tenancy"
)

func TestAddValidatesTenantPreservesInputAndRejectsCollisions(t *testing.T) {
	t.Parallel()

	event := tenancyEvent(t, nil)
	if _, err := Add(event, golibtenancy.TenantID{}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("zero tenant error = %v", err)
	}

	tenant := golibtenancy.MustTenantID("tenant-a")
	if _, err := Add(cloudevents.Event{}, tenant); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("zero event error = %v", err)
	}
	converted, err := Add(event, tenant)
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := event.Extension("tenantid"); ok {
		t.Fatal("input event gained tenantid")
	}
	if extracted, err := Extract(converted, true); err != nil || !extracted.Equal(tenant) {
		t.Fatalf("tenant = %v, %v", extracted, err)
	}
	if _, err := Add(converted, tenant); err != nil {
		t.Fatalf("equal tenant error = %v", err)
	}
	if _, err := Add(converted, golibtenancy.MustTenantID("tenant-b")); !errors.Is(err, ErrMetadataCollision) {
		t.Fatalf("collision error = %v", err)
	}
}

func TestExtractRequiresPresentTrustedValidStringMetadata(t *testing.T) {
	t.Parallel()

	if tenant, err := Extract(tenancyEvent(t, nil), true); !errors.Is(err, golibtenancy.ErrTenantMetadataMissing) || tenant.Valid() {
		t.Fatalf("missing tenant = %v, %v", tenant, err)
	}
	if tenant, err := Extract(tenancyEvent(t, map[string]cloudevents.Attribute{
		"tenantid": cloudevents.NewBooleanAttribute(true),
	}), true); !errors.Is(err, ErrInvalidInput) || tenant.Valid() {
		t.Fatalf("non-string tenant = %v, %v", tenant, err)
	}

	valid, err := cloudevents.NewStringAttribute("tenant-a")
	if err != nil {
		t.Fatal(err)
	}
	event := tenancyEvent(t, map[string]cloudevents.Attribute{"tenantid": valid})
	if tenant, err := Extract(event, false); !errors.Is(err, ErrUntrustedMetadata) || tenant.Valid() {
		t.Fatalf("untrusted tenant = %v, %v", tenant, err)
	}

	invalid, err := cloudevents.NewStringAttribute("bad?tenant")
	if err != nil {
		t.Fatal(err)
	}
	if tenant, err := Extract(tenancyEvent(t, map[string]cloudevents.Attribute{"tenantid": invalid}), true); !errors.Is(err, ErrInvalidInput) || tenant.Valid() {
		t.Fatalf("invalid tenant = %v, %v", tenant, err)
	}
	if tenant, err := Extract(tenancyEvent(t, map[string]cloudevents.Attribute{"tenantid": invalid}), false); !errors.Is(err, ErrUntrustedMetadata) || tenant.Valid() {
		t.Fatalf("untrusted invalid tenant = %v, %v", tenant, err)
	}
}

func tenancyEvent(t *testing.T, extensions map[string]cloudevents.Attribute) cloudevents.Event {
	t.Helper()
	event, err := cloudevents.NewEvent(cloudevents.Attributes{
		ID: "event-1", Source: "/source", Type: "example.created", Extensions: extensions,
	}, cloudevents.NewBinaryData([]byte("body")))
	if err != nil {
		t.Fatal(err)
	}
	return event
}
