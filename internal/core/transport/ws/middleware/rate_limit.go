package core_ws_middleware

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	core_redis_rate_limit "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/rate_limit"
)

type WSRateLimiter struct {
	limiter core_redis_rate_limit.RateLimiterMethods
}

type Result core_redis_rate_limit.Result

func NewWSRateLimiter(pool core_redis_pool.Pool, limit int64, window time.Duration) *WSRateLimiter {
	limiter := core_redis_rate_limit.New(pool, limit, window)

	return &WSRateLimiter{
		limiter: limiter,
	}
}

func (l *WSRateLimiter) Allow(ctx context.Context, eventType domain.RealtimeEventType, ip net.IP) (Result, error) {
	if l == nil || l.limiter == nil {
		return Result{}, fmt.Errorf("websocket auth rate limiter is not configured")
	}

	if eventType == domain.RealtimeEventNil {
		return Result{}, fmt.Errorf("`eventType` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if ip == nil {
		return Result{}, fmt.Errorf("`ip` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	scope := string(eventType)

	result, err := l.limiter.Allow(ctx, scope, ip)
	if err != nil {
		return Result{}, fmt.Errorf("limiter allow: %w", err)
	}

	return Result(result), nil
}
