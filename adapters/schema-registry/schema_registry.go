// Package schemaregistry provides an opt-in, bounded registry-backed JSON
// Schema validator for CloudEvents. Event decoding never performs lookup.
package schemaregistry

import (
	"context"
	"errors"
	"fmt"
	"maps"
	"time"

	"github.com/faustbrian/go-cloudevents"
	"github.com/faustbrian/go-cloudevents/internal/adapter"
	golibregistry "github.com/faustbrian/go-schema-registry"
	registryjsonschema "github.com/faustbrian/go-schema-registry/formats/jsonschema"
)

var (
	// ErrSchemaViolation reports that a resolved JSON Schema rejects a payload.
	// It aliases the CloudEvents facade sentinel for errors.Is compatibility.
	ErrSchemaViolation = cloudevents.ErrSchemaViolation
	// ErrSchemaMapping reports invalid validator configuration, an untrusted
	// dataschema URI, or a resolved schema with an unsupported format. It aliases
	// the CloudEvents facade sentinel.
	ErrSchemaMapping = cloudevents.ErrSchemaMapping
)

// JSONSchemaConfig defines the bounded, caller-owned registry dependencies and
// trust policy used to construct a JSONSchemaValidator.
type JSONSchemaConfig struct {
	// Cache resolves mapped lookups and owns registry availability and identity
	// checks. It must be non-nil.
	Cache *golibregistry.ResolveCache
	// SchemaLookups maps each exact trusted CloudEvents dataschema URI to the
	// registry lookup it may resolve. Construction clones the map, so later
	// caller mutations do not change the validator's trust boundary.
	SchemaLookups map[string]golibregistry.Lookup
	// Adapter decodes and validates resolved JSON Schema payloads. The caller
	// owns its limits, lifecycle, and referenced-schema configuration.
	Adapter *registryjsonschema.Adapter
	// AvailabilityPolicy controls whether resolution fails closed, accepts
	// stale cache entries, uses cache only, or reports unavailability.
	AvailabilityPolicy golibregistry.AvailabilityPolicy
	// Timeout is the positive upper bound applied to each resolution and payload
	// validation operation, subject to an earlier parent-context deadline.
	Timeout time.Duration
}

// JSONSchemaValidator validates caller-mapped CloudEvents dataschema URIs
// through a bounded schema-registry cache. Its zero value is intentionally
// unconfigured. Values returned by NewJSONSchemaValidator are safe for
// concurrent validation when their Cache and Adapter are safe for concurrent
// use.
type JSONSchemaValidator struct {
	cache              *golibregistry.ResolveCache
	schemaLookups      map[string]golibregistry.Lookup
	adapter            *registryjsonschema.Adapter
	availabilityPolicy golibregistry.AvailabilityPolicy
	timeout            time.Duration
	configured         bool
}

// NewJSONSchemaValidator validates config and takes an immutable snapshot of
// SchemaLookups. It rejects missing dependencies, empty URI or lookup entries,
// unsupported availability policies, and non-positive timeouts with an error
// wrapping ErrSchemaMapping.
func NewJSONSchemaValidator(config JSONSchemaConfig) (JSONSchemaValidator, error) {
	if config.Cache == nil || len(config.SchemaLookups) == 0 || config.Adapter == nil || !validAvailabilityPolicy(config.AvailabilityPolicy) || config.Timeout <= 0 {
		return JSONSchemaValidator{}, fmt.Errorf("%w: validator configuration", ErrSchemaMapping)
	}
	for uri, lookup := range config.SchemaLookups {
		if uri == "" || lookup.Kind() == "" {
			return JSONSchemaValidator{}, fmt.Errorf("%w: dataschema mapping", ErrSchemaMapping)
		}
	}
	return JSONSchemaValidator{cache: config.Cache, schemaLookups: maps.Clone(config.SchemaLookups), adapter: config.Adapter, availabilityPolicy: config.AvailabilityPolicy, timeout: config.Timeout, configured: true}, nil
}

// Validate implements cloudevents.SchemaValidator. It resolves only an exact
// URI present in the constructor-owned allowlist, accepts only JSON-compatible
// content types and JSON Schema registry results, and bounds work by Timeout.
// Mapping failures wrap ErrSchemaMapping; payload violations wrap
// ErrSchemaViolation; resolution, cancellation, and other adapter failures are
// preserved for errors.Is inspection.
func (validator JSONSchemaValidator) Validate(ctx context.Context, uri, contentType string, data []byte) error {
	if ctx == nil {
		return cloudevents.ErrContextRequired
	}
	if !validator.configured {
		return ErrSchemaMapping
	}
	if !adapter.IsJSONContentType(contentType) {
		return fmt.Errorf("%w: non-JSON content type", ErrSchemaMapping)
	}
	lookup, found := validator.schemaLookups[uri]
	if !found {
		return fmt.Errorf("%w: dataschema URI", ErrSchemaMapping)
	}
	resolveCtx, cancel := context.WithTimeout(ctx, validator.timeout)
	defer cancel()
	resolution, err := validator.cache.Resolve(resolveCtx, lookup, validator.availabilityPolicy)
	if err != nil {
		return err
	}
	resolved := resolution.Result
	if resolved.Schema.Definition().Format != golibregistry.FormatJSONSchema {
		return fmt.Errorf("%w: registry format", ErrSchemaMapping)
	}
	var target any
	if err := validator.adapter.Decode(resolveCtx, resolved.Schema, data, &target); err != nil {
		if errors.Is(err, registryjsonschema.ErrPayloadInvalid) {
			return fmt.Errorf("%w: %w", ErrSchemaViolation, err)
		}
		return err
	}
	return nil
}

func validAvailabilityPolicy(policy golibregistry.AvailabilityPolicy) bool {
	switch policy {
	case golibregistry.FailClosed, golibregistry.AllowStale, golibregistry.CacheOnly, golibregistry.ReturnUnavailable:
		return true
	default:
		return false
	}
}
