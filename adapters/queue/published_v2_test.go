package queue

import (
	"errors"
	"testing"

	published "github.com/faustbrian/go-tenancy/v2"
)

func TestPublishedV2TenantErrorIdentity(t *testing.T) {
	if _, err := validatedTenant("tenant with spaces"); !errors.Is(err, ErrInvalidInput) || !errors.Is(err, published.ErrInvalidTenantID) {
		t.Fatalf("tenant validation lost public v2 error identity: %v", err)
	}
}
