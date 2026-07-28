package core_redis_cache

import (
	"context"
	"fmt"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
)

func (r *Repository) DelCache(ctx context.Context, key string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	redisKey := fmt.Sprintf("%s:%s:%s", core_redis_pool.CacheFolder, r.prefix, key)

	if err := r.pool.Del(ctx, redisKey); err != nil {
		return fmt.Errorf("redis del: %w", err)
	}

	return nil
}
