package core_redis_cache

import "context"

func (c *RedisCache) Delete(ctx context.Context, key string) error {
	if err := validateKey(key); err != nil {
		return err
	}

	return c.deleteRedisKeys(ctx, []string{c.redisKey(key)})
}

// DeleteMany deletes keys in a single Redis DEL command. Empty input is a no-op.
func (c *RedisCache) DeleteMany(ctx context.Context, keys ...string) error {
	if len(keys) == 0 {
		return nil
	}

	redisKeys := make([]string, 0, len(keys))

	for _, key := range keys {
		if err := validateKey(key); err != nil {
			return err
		}

		redisKeys = append(redisKeys, c.redisKey(key))
	}

	return c.deleteRedisKeys(ctx, redisKeys)
}
