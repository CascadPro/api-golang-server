package core_redis_rate_limit

import (
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	core_rate_limit "github.com/CascadePro/api-golang-server/internal/core/rate_limit"
)

type RateLimiter struct {
	pool core_redis_pool.Pool
}

func New(pool core_redis_pool.Pool) *RateLimiter {
	return &RateLimiter{pool: pool}
}

var _ core_rate_limit.Limiter = (*RateLimiter)(nil)
