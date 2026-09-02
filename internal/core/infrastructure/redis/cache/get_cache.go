package core_redis_cache

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
)

func (r *Repository) GetCache(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	redisKey := fmt.Sprintf("%s:%s:%s", core_redis_pool.CacheFolder, r.prefix, key)

	dest, err := r.pool.Get(ctx, redisKey)
	if err != nil {
		if errors.Is(err, core_redis_pool.ErrNoValue) {
			return "", fmt.Errorf("get cache: %w", core_errors.ErrNotFound)
		}

		return "", fmt.Errorf("redis get: %w", err)
	}

	return dest, nil
}
