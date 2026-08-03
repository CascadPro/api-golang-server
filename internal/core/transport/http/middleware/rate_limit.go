package core_http_middleware

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type RateLimitConfig struct {
	// Максимальное количество запросов
	Limit int64
	// Временной интервал
	Window time.Duration
}

var luaScript string = `
	local key = KEYS[1]
	local limit = tonumber(ARGV[1])
	local window = tonumber(ARGV[2])

	local current = redis.call('GET', key)

	if current == false then
		redis.call('SET', key, 1)
		redis.call('EXPIRE', key, window)
		return {1, limit - 1, window}
	end

	local count = tonumber(current)

	if count >= limit then
		local ttl = redis.call('TTL', key)
		return {0, 0, ttl}
	end

	local new_count = redis.call('INCR', key)
	local ttl = redis.call('TTL', key)

	if ttl == -1 then
		redis.call('EXPIRE', key, window)
		ttl = window
	end

	return {1, limit - new_count, ttl}
`

func NewRateLimitConfig(limit int64, window time.Duration) *RateLimitConfig {
	return &RateLimitConfig{
		Limit:  limit,
		Window: window,
	}
}

func (cfg *RateLimitConfig) Middleware() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewResponseHandler(log, rw)

			clientIP, err := core_context.IPFromContext(ctx)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to get client ip address")
				return
			}

			rdb, err := core_redis_pool.New(r.Context(), core_redis_pool.NewConfigMust())
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to init redis connection pool")

				return
			}
			defer rdb.Close()

			key := fmt.Sprintf("%s:%s:%s", core_redis_pool.RateLimitFolder, r.URL.String(), clientIP)

			result, err := rdb.Eval(
				r.Context(),
				luaScript,
				[]string{key},
				cfg.Limit,
				int64(cfg.Window.Seconds()),
			)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to execute redis query")

				return
			}

			resultSlice := result.([]any)
			allowed := resultSlice[0].(int64) == 1
			remaining := resultSlice[1].(int64)
			ttl := resultSlice[2].(int64)

			resetTime := time.Now().Add(time.Duration(ttl) * time.Second).Unix()

			rw.Header().Add("X-RateLimit-Limit", strconv.FormatInt(cfg.Limit, 10))
			rw.Header().Add("X-RateLimit-Remaining", strconv.FormatInt(remaining, 10))
			rw.Header().Add("X-RateLimit-Reset", strconv.FormatInt(resetTime, 10))

			if !allowed {
				msg := fmt.Sprintf("try again after %v seconds", ttl)
				responseHandler.ErrorResponse(core_errors.ErrTooManyRequests, msg)

				return
			}

			next.ServeHTTP(rw, r)
		})
	}
}
