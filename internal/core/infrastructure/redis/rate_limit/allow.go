package core_redis_rate_limit

import (
	"context"
	"fmt"
	"net"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
)

func (l *RateLimiter) Allow(ctx context.Context, scope string, ip net.IP) (Result, error) {
	if ip == nil {
		return Result{}, fmt.Errorf("rate limit client ip is nil: %w", core_errors.ErrInvalidArgument)
	}

	key := fmt.Sprintf(
		"%s:%s:%s",
		core_redis_pool.RateLimitFolder,
		scope,
		ip.String(),
	)

	result, err := l.pool.Eval(
		ctx,
		luaScript,
		[]string{key},
		l.limit,
		int64(l.window.Seconds()),
	)
	if err != nil {
		return Result{}, fmt.Errorf("execute rate limit script: %w", err)
	}

	values, ok := result.([]any)
	if !ok || len(values) != 3 {
		return Result{}, fmt.Errorf("invalid rate limit redis response")
	}

	allowed, ok := values[0].(int64)
	if !ok {
		return Result{}, fmt.Errorf("invalid rate limit allowed value")
	}

	remaining, ok := values[1].(int64)
	if !ok {
		return Result{}, fmt.Errorf("invalid rate limit remaining value")
	}

	ttl, ok := values[2].(int64)
	if !ok {
		return Result{}, fmt.Errorf("invalid rate limit ttl value")
	}

	if ttl < 0 {
		ttl = 0
	}

	return Result{
		Allowed:    allowed == 1,
		Remaining:  remaining,
		RetryAfter: time.Duration(ttl) * time.Second,
	}, nil
}
