package core_redis_cache

import (
	"context"
	"fmt"
)

// Invalidate deletes every key under this cache namespace. It uses SCAN rather
// than KEYS to avoid blocking Redis while the namespace is large. Entries may
// be added concurrently during invalidation; callers should coordinate writes
// if they require a strict point-in-time flush.
func (c *RedisCache) Invalidate(ctx context.Context) error {
	pattern := c.redisKey("*")

	var cursor uint64
	var allKeys []string

	for {
		opCtx, cancel, err := c.operationContext(ctx)
		if err != nil {
			return err
		}

		keys, next, scanErr := c.pool.Scan(opCtx, cursor, pattern, defaultScanCount)

		cancel()

		if scanErr != nil {
			return fmt.Errorf("scan cache namespace %q: %w", c.namespace, scanErr)
		}

		allKeys = append(allKeys, keys...)

		cursor = next
		if cursor == 0 {
			break
		}
	}

	// Do not delete while SCAN is still running: mutating the keyspace during
	// iteration can cause Redis to skip entries while its hash table rehashes.
	for start := 0; start < len(allKeys); start += deleteBatchSize {
		end := min(start+deleteBatchSize, len(allKeys))
		if err := c.deleteRedisKeys(ctx, allKeys[start:end]); err != nil {
			return fmt.Errorf("invalidate cache namespace %q: %w", c.namespace, err)
		}
	}

	return nil
}
