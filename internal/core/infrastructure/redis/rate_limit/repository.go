package core_redis_rate_limit

import (
	"context"
	"net"
	"time"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
)

const (
	DefaultLimit  = int64(100)
	DefaultWindow = time.Minute
)

type RateLimiter struct {
	pool   core_redis_pool.Pool
	limit  int64
	window time.Duration
}

type RateLimiterMethods interface {
	Allow(ctx context.Context, scope string, ip net.IP) (Result, error)
}

func New(pool core_redis_pool.Pool, limit int64, window time.Duration) *RateLimiter {
	if limit <= 0 {
		limit = DefaultLimit
	}
	if window <= 0 {
		window = DefaultWindow
	}

	return &RateLimiter{
		pool:   pool,
		limit:  limit,
		window: window,
	}
}

func NewDefault(pool core_redis_pool.Pool) *RateLimiter {
	return &RateLimiter{
		pool:   pool,
		limit:  DefaultLimit,
		window: DefaultWindow,
	}
}
