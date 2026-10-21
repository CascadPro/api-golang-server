package core_http_middleware

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	core_redis_rate_limit "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/rate_limit"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_rate_limit "github.com/CascadePro/api-golang-server/internal/core/rate_limit"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type RateLimitOptions struct {
	FailureMode core_rate_limit.FailureMode
	Observer    core_rate_limit.Observer
}

type HTTPRateLimiter struct {
	limiter core_rate_limit.Limiter
	policy  core_rate_limit.Policy
	options RateLimitOptions
	initErr error
}

type HTTPRateLimitMiddleware interface {
	Middleware() Middleware
}

// NewHTTPRateLimiter is retained for source compatibility. New app code should inject
// a shared Limiter through NewHTTPRateLimiterWithOptions.
func NewHTTPRateLimiter(pool core_redis_pool.Pool, limit int64, window time.Duration) *HTTPRateLimiter {
	policy := core_rate_limit.Policy{Name: "http.custom", Limit: limit, Window: window}
	rl, err := NewHTTPRateLimiterWithOptions(core_redis_rate_limit.New(pool), policy, RateLimitOptions{
		FailureMode: core_rate_limit.FailureClosed,
		Observer:    core_rate_limit.NoopObserver{},
	})
	if err != nil {
		return &HTTPRateLimiter{policy: policy, options: RateLimitOptions{FailureMode: core_rate_limit.FailureClosed}, initErr: err}
	}
	return rl
}

func NewHTTPRateLimiterWithOptions(limiter core_rate_limit.Limiter, policy core_rate_limit.Policy, options RateLimitOptions) (*HTTPRateLimiter, error) {
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
	return &HTTPRateLimiter{limiter: limiter, policy: policy, options: options}, nil
}

func (rl *HTTPRateLimiter) Middleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			locale := core_context.Locale(ctx)
			responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

			if rl == nil || rl.initErr != nil || rl.limiter == nil {
				err := fmt.Errorf("HTTP rate limiter is not configured")
				if rl != nil && rl.initErr != nil {
					err = rl.initErr
				}
				responseHandler.ErrorResponse(err, "rate limiter configuration error")
				return
			}

			clientIP, err := core_context.ClientIP(ctx)
			if err != nil || clientIP == nil {
				if err == nil {
					err = fmt.Errorf("client IP is nil")
				}
				rl.options.Observer.Observe(core_rate_limit.Observation{Policy: rl.policy.Name, Outcome: core_rate_limit.OutcomeError})
				responseHandler.ErrorResponse(err, "failed to get client IP address")
				return
			}

			started := time.Now()
			key := httpRateLimitKey(r.Method, r.URL.EscapedPath(), clientIP)
			decision, err := rl.limiter.Allow(ctx, key, rl.policy)
			elapsed := time.Since(started)
			if err != nil {
				rl.options.Observer.Observe(core_rate_limit.Observation{Policy: rl.policy.Name, Outcome: core_rate_limit.OutcomeError, Duration: elapsed})
				if rl.options.FailureMode == core_rate_limit.FailureOpen {
					rl.options.Observer.Observe(core_rate_limit.Observation{Policy: rl.policy.Name, Outcome: core_rate_limit.OutcomeBypass})
					next.ServeHTTP(rw, r)
					return
				}
				responseHandler.ErrorResponse(err, "failed to check rate limit")
				return
			}

			outcome := core_rate_limit.OutcomeAllowed
			if !decision.Allowed {
				outcome = core_rate_limit.OutcomeLimited
			}
			rl.options.Observer.Observe(core_rate_limit.Observation{Policy: rl.policy.Name, Outcome: outcome, Duration: elapsed})
			setRateLimitHeaders(rw.Header(), decision)
			if !decision.Allowed {
				rw.Header().Set("Retry-After", strconv.FormatInt(ceilSeconds(decision.RetryAfter), 10))
				responseHandler.ErrorResponse(core_errors.ErrTooManyRequests, "request rate limit exceeded")
				return
			}
			next.ServeHTTP(rw, r)
		})
	}
}

func httpRateLimitKey(method, path string, ip net.IP) string {
	return "http:" + method + ":" + path + ":" + ip.String()
}

func setRateLimitHeaders(header http.Header, decision core_rate_limit.Decision) {
	header.Set("X-RateLimit-Limit", strconv.FormatInt(decision.Limit, 10))
	header.Set("X-RateLimit-Remaining", strconv.FormatInt(decision.Remaining, 10))
	header.Set("X-RateLimit-Reset", strconv.FormatInt(time.Now().Add(decision.RetryAfter).Unix(), 10))
}

func ceilSeconds(d time.Duration) int64 {
	if d <= 0 {
		return 0
	}
	return int64((d + time.Second - 1) / time.Second)
}
