package core_redis_cache

import (
	"context"
	"encoding/json"
	"strings"
	"sync"
	"time"
)

const (
	defaultScanCount = int64(500)
	deleteBatchSize  = 500
)

// RedisCommands is the narrow Redis dependency required by this cache.
// The full Redis pool remains available to sessions, Pub/Sub, rate limiting
// and other infrastructure.
type RedisCommands interface {
	Get(ctx context.Context, key string) (string, error)
	MGet(ctx context.Context, keys ...string) ([]any, error)

	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	SetMany(ctx context.Context, values map[string]any, expiration time.Duration) error

	Del(ctx context.Context, keys ...string) error

	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)

	OpTimeout() time.Duration
}

type loadCall struct {
	done chan struct{}
	data json.RawMessage
	err  error
}

type RedisCache struct {
	pool       RedisCommands
	namespace  string
	defaultTTL time.Duration

	flightsMu sync.Mutex
	flights   map[string]*loadCall
}

// NewRedisCache creates a cache isolated by namespace. A zero TTL on Set uses
// defaultTTL; when both are zero Redis stores the key without expiration.
func NewRedisCache(pool RedisCommands, namespace string, defaultTTL time.Duration) (*RedisCache, error) {
	if isNilValue(pool) {
		return nil, ErrNilPool
	}

	namespace = strings.TrimSpace(namespace)
	if !validNamespace(namespace) {
		return nil, ErrInvalidNamespace
	}

	if defaultTTL < 0 {
		return nil, ErrInvalidTTL
	}

	return &RedisCache{
		pool:       pool,
		namespace:  namespace,
		defaultTTL: defaultTTL,
		flights:    make(map[string]*loadCall),
	}, nil
}

var (
	_ Cache                = (*RedisCache)(nil)
	_ BatchCache           = (*RedisCache)(nil)
	_ NamespaceInvalidator = (*RedisCache)(nil)
	_ LoaderCache          = (*RedisCache)(nil)
)
