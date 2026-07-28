package core_redis_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisFolderName string

const (
	SessionFolder   = RedisFolderName("cascade__session")
	CacheFolder     = RedisFolderName("cascade__cache")
	RateLimitFolder = RedisFolderName("cascade__rate_limit")
)

type Pool interface {
	Get(ctx context.Context, key string) (string, error)
	GetKeys(ctx context.Context, cursor uint64, match string, count int64) ([]string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Del(ctx context.Context, keys ...string) error
	Eval(ctx context.Context, script string, keys []string, args ...any) (any, error)

	HGet(ctx context.Context, key, field string) (string, error)
	HGetAll(ctx context.Context, key string) *redis.MapStringStringCmd
	HSet(ctx context.Context, key string, values ...any) error
	HDel(ctx context.Context, key string, fields ...string) error

	Pipeline() redis.Pipeliner
	TxPipeline() redis.Pipeliner

	Close()

	OpTimeout() time.Duration
}

type ConnectionPool struct {
	client  *redis.Client
	timeout time.Duration
}

func NewConnectionPool(ctx context.Context, cfg Config) (*ConnectionPool, error) {
	addr := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: cfg.Password,
		DB:       cfg.Database,

		PoolSize:     cfg.PoolSize,
		MinIdleConns: cfg.MinIdleConns,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  5 * time.Second,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: %w", err)
	}

	return &ConnectionPool{
		client:  rdb,
		timeout: cfg.Timeout,
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.timeout
}
