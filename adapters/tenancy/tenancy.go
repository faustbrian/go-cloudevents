// Package tenancy maps trusted Golib tenant routing identity to and from a
// CloudEvents extension. It does not authenticate or authorize the tenant.
package tenancy

import (
	"fmt"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	golibtenancy "github.com/faustbrian/go-tenancy"
)

var (
	// ErrInvalidInput identifies an invalid tenant or malformed tenant
	// extension. Callers can match it with errors.Is.
	ErrInvalidInput = cloudevents.ErrInvalidAdapterInput
	// ErrMetadataCollision identifies an event whose existing tenant extension
	// conflicts with the tenant being added. The input event remains unchanged.
	ErrMetadataCollision = cloudevents.ErrMetadataCollision
	// ErrUntrustedMetadata identifies an attempt to adopt a tenant extension
	// before the caller has established trust in the event's transport boundary.
	ErrUntrustedMetadata = cloudevents.ErrUntrustedMetadata
)

// Add returns a new event with one validated tenant routing extension. The
// zero TenantID or zero Event returns ErrInvalidInput. An equal existing value
// is accepted; a conflicting value returns ErrMetadataCollision. The caller
// retains ownership of event and tenant, neither of which is mutated.
func Add(event cloudevents.Event, tenant golibtenancy.TenantID) (cloudevents.Event, error) {
	if !tenant.Valid() {
		return cloudevents.Event{}, fmt.Errorf("%w: tenant", ErrInvalidInput)
	}
	return adapter.AddStringExtensions(event, map[string]string{adapter.TenantIDExtension: tenant.Value()})
}

// Extract adopts the tenant extension only after an explicit trust decision.
// Missing metadata or a malformed extension kind returns ErrInvalidInput.
// trusted=false returns ErrUntrustedMetadata without validating or adopting a
// string value; a malformed trusted tenant returns ErrInvalidInput. The zero
// TenantID is returned on every failure. The input event remains caller-owned
// and is not mutated.
func Extract(event cloudevents.Event, trusted bool) (golibtenancy.TenantID, error) {
	value, present, err := adapter.StringExtension(event, adapter.TenantIDExtension)
	if err != nil {
		return golibtenancy.TenantID{}, err
	}
	if !present {
		return golibtenancy.TenantID{}, fmt.Errorf("%w: tenantid: %w", ErrInvalidInput, golibtenancy.ErrTenantMetadataMissing)
	}
	if !trusted {
		return golibtenancy.TenantID{}, ErrUntrustedMetadata
	}
	tenant, err := golibtenancy.ParseTenantID(value)
	if err != nil {
		return golibtenancy.TenantID{}, fmt.Errorf("%w: tenantid: %w", ErrInvalidInput, err)
	}
	return tenant, nil
}
