package core_redis_cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

func (c *RedisCache) Set(ctx context.Context, key string, value any, ttl time.Duration) error {
	if err := validateKey(key); err != nil {
		return err
	}
	ttl, err := c.effectiveTTL(ttl)
	if err != nil {
		return err
	}

	encoded, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("encode cached JSON: %w", err)
	}

	opCtx, cancel, err := c.operationContext(ctx)
	if err != nil {
		return err
	}
	defer cancel()

	if err := c.pool.Set(opCtx, c.redisKey(key), encoded, ttl); err != nil {
		return fmt.Errorf("cache set in Redis: %w", err)
	}

	return nil
}

// SetMany serializes all values before writing anything, then stores them with
// a single pipelined Redis operation. An empty input is a no-op.
func (c *RedisCache) SetMany(ctx context.Context, values map[string]any, ttl time.Duration) error {
	if len(values) == 0 {
		return nil
	}
	ttl, err := c.effectiveTTL(ttl)
	if err != nil {
		return err
	}

	encoded := make(map[string]any, len(values))

	for key, value := range values {
		if err := validateKey(key); err != nil {
			return fmt.Errorf("cache set many key %q: %w", key, err)
		}

		data, err := json.Marshal(value)
		if err != nil {
			return fmt.Errorf("encode cached JSON for %q: %w", key, err)
		}

		encoded[c.redisKey(key)] = data
	}

	opCtx, cancel, err := c.operationContext(ctx)
	if err != nil {
		return err
	}
	defer cancel()

	if err := c.pool.SetMany(opCtx, encoded, ttl); err != nil {
		return fmt.Errorf("cache set many in Redis: %w", err)
	}

	return nil
}
