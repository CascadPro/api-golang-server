package core_redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
)

func (c *RedisCache) Get(ctx context.Context, key string, dest any) error {
	if err := validateKey(key); err != nil {
		return err
	}
	if isNilValue(dest) {
		return fmt.Errorf("cache get destination: %w", ErrNilDestination)
	}

	opCtx, cancel, err := c.operationContext(ctx)
	if err != nil {
		return err
	}
	defer cancel()

	value, err := c.pool.Get(opCtx, c.redisKey(key))
	if err != nil {
		if errors.Is(err, core_redis_pool.ErrNoValue) {
			return fmt.Errorf("cache get: %w", core_errors.ErrNotFound)
		}

		return fmt.Errorf("cache get from Redis: %w", err)
	}

	if err := json.Unmarshal([]byte(value), dest); err != nil {
		return fmt.Errorf("decode cached JSON: %w", err)
	}

	return nil
}

// GetMany performs one Redis MGET and decodes each present value into its
// corresponding destination. The returned map contains per-key failures;
// absent keys map to core_errors.ErrNotFound. An empty input is a no-op.
func (c *RedisCache) GetMany(ctx context.Context, destinations map[string]any) (map[string]error, error) {
	failures := make(map[string]error)
	if len(destinations) == 0 {
		return failures, nil
	}

	keys := make([]string, 0, len(destinations))
	originalKeys := make([]string, 0, len(destinations))

	for key, dest := range destinations {
		if err := validateKey(key); err != nil {
			return nil, fmt.Errorf("cache get many key %q: %w", key, err)
		}
		if isNilValue(dest) {
			return nil, fmt.Errorf("cache get many destination for %q: %w", key, ErrNilDestination)
		}

		originalKeys = append(originalKeys, key)
		keys = append(keys, c.redisKey(key))
	}

	opCtx, cancel, err := c.operationContext(ctx)
	if err != nil {
		return nil, err
	}
	defer cancel()

	values, err := c.pool.MGet(opCtx, keys...)
	if err != nil {
		return nil, fmt.Errorf("cache MGET from Redis: %w", err)
	}
	if len(values) != len(originalKeys) {
		return nil, fmt.Errorf("cache MGET returned %d values for %d keys", len(values), len(originalKeys))
	}

	for i, raw := range values {
		key := originalKeys[i]
		if raw == nil {
			failures[key] = fmt.Errorf("cache get many: %w", core_errors.ErrNotFound)
			continue
		}

		var encoded []byte

		switch value := raw.(type) {
		case string:
			encoded = []byte(value)
		case []byte:
			encoded = value
		default:
			failures[key] = fmt.Errorf("cache MGET returned unsupported value type %T", raw)
			continue
		}

		if err := json.Unmarshal(encoded, destinations[key]); err != nil {
			failures[key] = fmt.Errorf("decode cached JSON for %q: %w", key, err)
		}
	}

	return failures, nil
}
