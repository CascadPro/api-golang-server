package core_redis_pool

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
)

var (
	ErrCacheMiss = errors.New("cache miss")
	ErrRedisDown = errors.New("redis is down")
)

func MapErrors(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		return ErrCacheMiss
	}

	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return fmt.Errorf("network failure: %w: %w", ErrRedisDown, err)
	}

	return fmt.Errorf("redis operation failed: %w", err)
}

func (p *ConnectionPool) Get(ctx context.Context, key string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	val, err := p.client.Get(ctx, key).Result()
	return val, MapErrors(err)
}

func (p *ConnectionPool) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	ctx, cancel := context.WithTimeout(ctx, p.timeout)
	defer cancel()

	err := p.client.Set(ctx, key, value, expiration).Err()
	return MapErrors(err)
}

func (p *ConnectionPool) Close() {
	_ = p.client.Close()
}
