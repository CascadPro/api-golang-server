package core_ipinfo_pool

import (
	"context"
	"encoding/json"
	"fmt"

	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/repository/redis/cache"
	"github.com/ipinfo/go/v2/ipinfo"
	"github.com/ipinfo/go/v2/ipinfo/cache"
)

type CacheEngine struct {
	ctx  context.Context
	repo core_redis_cache.RepositoryMethods

	cache.Interface
}

func NewCacheEngine(ctx context.Context, repo core_redis_cache.RepositoryMethods) *CacheEngine {
	return &CacheEngine{
		ctx:  ctx,
		repo: repo,
	}
}

func (ce *CacheEngine) Get(key string) (any, error) {
	jsonString, err := ce.repo.GetCache(ce.ctx, key)
	if err != nil {
		return nil, fmt.Errorf("engine cache: %w", err)
	}

	var dest ipinfo.Core
	if err := json.Unmarshal([]byte(jsonString), &dest); err != nil {
		return nil, fmt.Errorf("engine cache: %w", err)
	}

	return &dest, nil
}

func (ce *CacheEngine) Set(key string, value any) error {
	return ce.repo.SetCache(ce.ctx, key, value)
}
