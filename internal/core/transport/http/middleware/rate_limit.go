package core_http_middleware

import (
	"net/http"
	"strconv"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	core_redis_rate_limit "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/rate_limit"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type HTTPRateLimiter struct {
	repository core_redis_rate_limit.RateLimiterMethods
}

type HTTPRateLimitMiddleware interface {
	Middleware() Middleware
}

func NewHTTPRateLimiter(pool core_redis_pool.Pool, limit int64, window time.Duration) *HTTPRateLimiter {
	repository := core_redis_rate_limit.New(pool, limit, window)

	return &HTTPRateLimiter{
		repository: repository,
	}
}

func (rl *HTTPRateLimiter) Middleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			locale := core_context.Locale(ctx)

			responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

			clientIP, err := core_context.ClientIP(ctx)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to get client ip address")
				return
			}

			result, err := rl.repository.Allow(ctx, r.URL.String(), clientIP)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to do repository execution")
				return
			}

			remaining := strconv.FormatInt(result.Remaining, 10)
			reset := strconv.FormatInt(time.Now().Add(result.RetryAfter).Unix(), 10)
			retryAfter := strconv.FormatInt(result.RetryAfter.Milliseconds(), 10)

			rw.Header().Add("X-RateLimit-Remaining", remaining)
			rw.Header().Add("X-RateLimit-Reset", reset)
			rw.Header().Add("X-RateLimit-Retry-After", retryAfter)

			if !result.Allowed {
				responseHandler.ErrorResponse(core_errors.ErrTooManyRequests, "requests limit is overrated")
				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
