package core_redis_cache

import (
	"context"
	"fmt"
	"reflect"
	"strings"
	"time"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
)

func (c *RedisCache) deleteRedisKeys(ctx context.Context, keys []string) error {
	if len(keys) == 0 {
		return nil
	}

	opCtx, cancel, err := c.operationContext(ctx)
	if err != nil {
		return err
	}
	defer cancel()

	if err := c.pool.Del(opCtx, keys...); err != nil {
		return fmt.Errorf("cache delete in Redis: %w", err)
	}

	return nil
}

func (c *RedisCache) effectiveTTL(ttl time.Duration) (time.Duration, error) {
	if ttl < 0 {
		return 0, ErrInvalidTTL
	}
	if ttl == 0 {
		ttl = c.defaultTTL
	}

	return ttl, nil
}

func (c *RedisCache) redisKey(key string) string {
	return fmt.Sprintf("%s:%s:%s", core_redis_pool.CacheFolder, c.namespace, key)
}

func (c *RedisCache) operationContext(ctx context.Context) (context.Context, context.CancelFunc, error) {
	if ctx == nil {
		return nil, nil, ErrNilContext
	}

	timeout := c.pool.OpTimeout()
	if timeout <= 0 {
		return ctx, func() {}, nil
	}

	opCtx, cancel := context.WithTimeout(ctx, timeout)

	return opCtx, cancel, nil
}

func validateKey(key string) error {
	if strings.TrimSpace(key) == "" {
		return ErrInvalidKey
	}

	return nil
}

func validNamespace(namespace string) bool {
	if namespace == "" {
		return false
	}

	for _, r := range namespace {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_' {
			continue
		}

		return false
	}

	return true
}

func isNilValue(value any) bool {
	if value == nil {
		return true
	}

	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Pointer, reflect.Slice:
		return v.IsNil()
	default:
		return false
	}
}
