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
	ErrNoValue   = errors.New("redis no value")
	ErrRedisDown = errors.New("redis is down")
	ErrUnknown   = errors.New("unknown redis error")
)

func MapErrors(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		return ErrNoValue
	}

	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return fmt.Errorf("network failure: %w: %w", ErrRedisDown, err)
	}

	return fmt.Errorf("redis operation failed: %v: %w", err, ErrUnknown)
}

func (p *ConnectionPool) Get(ctx context.Context, key string) (string, error) {
	val, err := p.client.Get(ctx, key).Result()
	return val, MapErrors(err)
}

func (p *ConnectionPool) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return MapErrors(p.client.Set(ctx, key, value, expiration).Err())
}

func (p *ConnectionPool) Del(ctx context.Context, keys ...string) error {
	return MapErrors(p.client.Del(ctx, keys...).Err())
}

func (p *ConnectionPool) Eval(ctx context.Context, script string, keys []string, args ...any) (any, error) {
	res, err := p.client.Eval(ctx, script, keys, args...).Result()
	return res, MapErrors(err)
}

func (p *ConnectionPool) GetKeys(ctx context.Context, cursor uint64, match string, count int64) ([]string, error) {
	var keys []string

	for {
		slice, c, err := p.client.Scan(ctx, cursor, match, count).Result()
		if err != nil {
			return nil, MapErrors(err)
		}
		cursor = c
		keys = append(keys, slice...)

		if cursor == 0 {
			break
		}
	}

	if len(keys) == 0 {
		return nil, MapErrors(redis.Nil)
	}

	return keys, nil
}

func (p *ConnectionPool) HGet(ctx context.Context, key, field string) (string, error) {
	return p.client.HGet(ctx, key, field).Result()
}

func (p *ConnectionPool) HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd {
	return p.client.HGetAll(ctx, key)
}

func (p *ConnectionPool) HSet(ctx context.Context, key string, values ...any) error {
	return MapErrors(p.client.HSet(ctx, key, values...).Err())
}

func (p *ConnectionPool) HDel(ctx context.Context, key string, fields ...string) error {
	return MapErrors(p.client.HDel(ctx, key, fields...).Err())
}

func (p *ConnectionPool) Pipeline() redis.Pipeliner {
	return p.client.Pipeline()
}

func (p *ConnectionPool) TxPipeline() redis.Pipeliner {
	return p.client.TxPipeline()
}

func (p *ConnectionPool) Close() {
	_ = p.client.Close()
}
