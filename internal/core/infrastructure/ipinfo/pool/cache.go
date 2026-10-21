package core_ipinfo_pool

import (
	"context"
	"fmt"

	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/cache"
	"github.com/ipinfo/go/v2/ipinfo"
)

// CacheEngine adapts the application's Cache contract to the interface expected by ipinfo-go.
// The upstream interface has no context parameter, so each operation gets a fresh context.
type CacheEngine struct {
	repository core_redis_cache.Cache
}

func NewCacheEngine(repo core_redis_cache.Cache) *CacheEngine {
	return &CacheEngine{repository: repo}
}

func (ce *CacheEngine) Get(key string) (any, error) {
	ctx := context.Background()

	var dest ipinfo.Core
	if err := ce.repository.Get(ctx, key, &dest); err != nil {
		return nil, fmt.Errorf("IPInfo cache get: %w", err)
	}

	return &dest, nil
}

func (ce *CacheEngine) Set(key string, value any) error {
	ctx := context.Background()
	if err := ce.repository.Set(ctx, key, value, 0); err != nil {
		return fmt.Errorf("IPInfo cache set: %w", err)
	}

	return nil
}
