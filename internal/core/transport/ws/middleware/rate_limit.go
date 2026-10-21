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
	core_rate_limit "github.com/CascadePro/api-golang-server/internal/core/rate_limit"
)

type WSRateLimiter struct {
	limiter core_rate_limit.Limiter
	policy  core_rate_limit.Policy
	options core_rate_limit.Options
	initErr error
}

type Result = core_rate_limit.Decision

// NewWSRateLimiter is kept for compatibility with current call sites.
func NewWSRateLimiter(pool core_redis_pool.Pool, limit int64, window time.Duration) *WSRateLimiter {
	l, err := NewWSRateLimiterWithOptions(core_redis_rate_limit.New(pool),
		core_rate_limit.Policy{Name: "ws.custom", Limit: limit, Window: window},
		core_rate_limit.Options{FailureMode: core_rate_limit.FailureClosed})
	if err != nil {
		return &WSRateLimiter{initErr: err}
	}
	return l
}

func NewWSRateLimiterWithOptions(limiter core_rate_limit.Limiter, policy core_rate_limit.Policy, options core_rate_limit.Options) (*WSRateLimiter, error) {
	if limiter == nil {
		return nil, fmt.Errorf("rate limiter is nil")
	}
	if err := policy.Validate(); err != nil {
		return nil, err
	}
	if options.FailureMode == "" {
		options.FailureMode = core_rate_limit.FailureClosed
	}
	if err := options.FailureMode.Validate(); err != nil {
		return nil, err
	}
	if options.Observer == nil {
		options.Observer = core_rate_limit.NoopObserver{}
	}
	return &WSRateLimiter{limiter: limiter, policy: policy, options: options}, nil
}

func (l *WSRateLimiter) Allow(ctx context.Context, eventType domain.RealtimeEventType, ip net.IP) (Result, error) {
	if l == nil {
		return Result{}, fmt.Errorf("websocket rate limiter is not configured")
	}
	if l.initErr != nil {
		return Result{}, fmt.Errorf("invalid websocket rate limiter configuration: %w", l.initErr)
	}
	if l.limiter == nil {
		return Result{}, fmt.Errorf("websocket rate limiter is not configured")
	}
	if eventType == domain.RealtimeEventNil {
		return Result{}, fmt.Errorf("event type is nil: %w", core_errors.ErrInvalidArgument)
	}
	if ip == nil {
		return Result{}, fmt.Errorf("client IP is nil: %w", core_errors.ErrInvalidArgument)
	}

	started := time.Now()
	key := "ws:" + string(eventType) + ":" + ip.String()
	result, err := l.limiter.Allow(ctx, key, l.policy)
	outcome := core_rate_limit.OutcomeAllowed
	if err != nil {
		outcome = core_rate_limit.OutcomeError
	} else if !result.Allowed {
		outcome = core_rate_limit.OutcomeLimited
	}
	l.options.Observer.Observe(core_rate_limit.Observation{Policy: l.policy.Name, Outcome: outcome, Duration: time.Since(started)})
	if err != nil {
		if l.options.FailureMode == core_rate_limit.FailureOpen {
			return core_rate_limit.Decision{Allowed: true, Limit: l.policy.Limit}, nil
		}
		return Result{}, fmt.Errorf("check websocket rate limit: %w", err)
	}
	return result, nil
}
