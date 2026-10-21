package core_redis_rate_limit

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strconv"
	"time"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
	core_rate_limit "github.com/CascadePro/api-golang-server/internal/core/rate_limit"
)

func (l *RateLimiter) Allow(ctx context.Context, key string, policy core_rate_limit.Policy) (core_rate_limit.Decision, error) {
	if l == nil || l.pool == nil {
		return core_rate_limit.Decision{}, errors.New("rate limiter Redis pool is not configured")
	}

	if err := policy.Validate(); err != nil {
		return core_rate_limit.Decision{}, err
	}

	if key == "" {
		return core_rate_limit.Decision{}, core_rate_limit.ErrInvalidKey
	}

	sum := sha256.Sum256([]byte(key))
	redisKey := fmt.Sprintf("%s:%s", core_redis_pool.RateLimitFolder, hex.EncodeToString(sum[:]))

	raw, err := l.pool.Eval(ctx, luaScript, []string{redisKey}, policy.Limit, policy.Window.Milliseconds())
	if err != nil {
		return core_rate_limit.Decision{}, fmt.Errorf("execute rate limit script: %w", err)
	}

	values, ok := raw.([]any)
	if !ok || len(values) != 3 {
		return core_rate_limit.Decision{}, fmt.Errorf("invalid Redis rate limit response type %T", raw)
	}

	allowed, err := redisInt64(values[0])
	if err != nil {
		return core_rate_limit.Decision{}, fmt.Errorf("parse rate limit allowed value: %w", err)
	}

	remaining, err := redisInt64(values[1])
	if err != nil {
		return core_rate_limit.Decision{}, fmt.Errorf("parse rate limit remaining value: %w", err)
	}

	ttlMillis, err := redisInt64(values[2])
	if err != nil {
		return core_rate_limit.Decision{}, fmt.Errorf("parse rate limit TTL: %w", err)
	}

	if ttlMillis < 0 {
		ttlMillis = 0
	}

	if remaining < 0 {
		remaining = 0
	}

	return core_rate_limit.Decision{
		Allowed: allowed == 1, Limit: policy.Limit, Remaining: remaining,
		RetryAfter: time.Duration(ttlMillis) * time.Millisecond,
	}, nil
}

func redisInt64(value any) (int64, error) {
	switch v := value.(type) {
	case int64:
		return v, nil
	case int:
		return int64(v), nil
	case string:
		return strconv.ParseInt(v, 10, 64)
	case []byte:
		return strconv.ParseInt(string(v), 10, 64)
	default:
		return 0, fmt.Errorf("expected integer, got %T", value)
	}
}
