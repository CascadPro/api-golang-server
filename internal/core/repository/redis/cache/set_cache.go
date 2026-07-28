package core_redis_cache

import (
	"context"
	"encoding/json"
	"fmt"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
)

func (r *Repository) SetCache(ctx context.Context, key string, dest any) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	redisKey := fmt.Sprintf("%s:%s:%s", core_redis_pool.CacheFolder, r.prefix, key)

	jsonString, err := json.Marshal(dest)
	if err != nil {
		return fmt.Errorf("json encode: %w", err)
	}

	if err := r.pool.Set(ctx, redisKey, jsonString, r.expiration); err != nil {
		return fmt.Errorf("redis set: %w", err)
	}

	return nil
}
