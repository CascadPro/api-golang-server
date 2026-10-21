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
	ErrTimedOut  = errors.New("redis operation timed out")
	ErrCanceled  = errors.New("redis operation canceled")
	ErrUnknown   = errors.New("unknown redis error")
)

func MapErrors(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, redis.Nil) {
		return ErrNoValue
	}

	if errors.Is(err, context.DeadlineExceeded) {
		return fmt.Errorf("%v: %w", err, ErrTimedOut)
	}

	if errors.Is(err, context.Canceled) {
		return fmt.Errorf("%v: %w", err, ErrCanceled)
	}

	var netErr *net.OpError
	if errors.As(err, &netErr) {
		return fmt.Errorf("network failure: %v: %w", err, ErrRedisDown)
	}

	// Preserve both the generic Redis classification and the original error
	// so callers can still use errors.Is/errors.As.
	return fmt.Errorf("redis operation failed: %w: %w", ErrUnknown, err)
}

func (p *ConnectionPool) Get(ctx context.Context, key string) (string, error) {
	val, err := p.client.Get(ctx, key).Result()
	return val, MapErrors(err)
}

func (p *ConnectionPool) Set(ctx context.Context, key string, value any, expiration time.Duration) error {
	return MapErrors(p.client.Set(ctx, key, value, expiration).Err())
}

func (p *ConnectionPool) SetMany(ctx context.Context, values map[string]any, expiration time.Duration) error {
	pipe := p.TxPipeline()

	for key, value := range values {
		if err := pipe.Set(ctx, key, value, expiration).Err(); err != nil {
			return MapErrors(err)
		}
	}

	if _, err := pipe.Exec(ctx); err != nil {
		return MapErrors(fmt.Errorf("SetMany pipe: %w", err))
	}

	return nil
}

func (p *ConnectionPool) Del(ctx context.Context, keys ...string) error {
	return MapErrors(p.client.Del(ctx, keys...).Err())
}

func (p *ConnectionPool) Eval(ctx context.Context, script string, keys []string, args ...any) (any, error) {
	res, err := p.client.Eval(ctx, script, keys, args...).Result()
	return res, MapErrors(err)
}

func (p *ConnectionPool) Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error) {
	cmd := p.client.Scan(ctx, cursor, match, count)
	keys, cursor := cmd.Val()

	return keys, cursor, MapErrors(cmd.Err())
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

func (p *ConnectionPool) MGet(ctx context.Context, keys ...string) ([]any, error) {
	result, err := p.client.MGet(ctx, keys...).Result()
	return result, MapErrors(err)
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

func (p *ConnectionPool) Publish(ctx context.Context, channel string, message any) error {
	return MapErrors(p.client.Publish(ctx, channel, message).Err())
}

func (p *ConnectionPool) Subscribe(ctx context.Context, channels ...string) *redis.PubSub {
	return p.client.Subscribe(ctx, channels...)
}

func (p *ConnectionPool) Close() {
	_ = p.client.Close()
}
