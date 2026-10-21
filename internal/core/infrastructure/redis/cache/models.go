package core_redis_cache

import (
	"context"
	"time"
)

// Cache is the application-facing contract for storing JSON-serializable values.
// Implementations hide key construction, serialization and the backing store.
type Cache interface {
	Get(ctx context.Context, key string, dest any) error
	Set(ctx context.Context, key string, value any, ttl time.Duration) error
	Delete(ctx context.Context, key string) error
}

// BatchCache is an optional extension for callers that need bulk operations.
// Per-key GetMany failures (for example, a missing key or invalid JSON) are
// returned in the result map. A non-nil top-level error means the Redis request
// itself failed or the input could not be prepared.
type BatchCache interface {
	GetMany(ctx context.Context, destinations map[string]any) (map[string]error, error)
	SetMany(ctx context.Context, values map[string]any, ttl time.Duration) error
	DeleteMany(ctx context.Context, keys ...string) error
}

// NamespaceInvalidator is an optional extension for invalidating every entry
// owned by a single cache namespace.
type NamespaceInvalidator interface {
	Invalidate(ctx context.Context) error
}

// LoaderCache prevents duplicate concurrent cache-miss loads in one process.
// The loader must return a JSON-serializable value matching dest's type.
type LoaderCache interface {
	GetOrLoad(ctx context.Context, key string, dest any, ttl time.Duration, loader func(context.Context) (any, error)) error
}
