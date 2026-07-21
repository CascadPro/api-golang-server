package core_redis_pool

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Pool interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value any, expiration time.Duration) error
	Close()

	OpTimeout() time.Duration
}

type ConnectionPool struct {
	client  *redis.Client
	timeout time.Duration
}

func NewConnectionPool(ctx context.Context, config Config) (*ConnectionPool, error) {
	addr := fmt.Sprintf("%s:%s", config.Host, config.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: config.Password,
		DB:       config.Database,

		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,

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
		timeout: config.Timeout,
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.timeout
}
