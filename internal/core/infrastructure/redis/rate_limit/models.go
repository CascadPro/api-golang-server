package core_redis_rate_limit

import (
	"context"

	core_rate_limit "github.com/CascadePro/api-golang-server/internal/core/rate_limit"
)

type Result = core_rate_limit.Decision

type RateLimiterMethods interface {
	Allow(ctx context.Context, key string, policy core_rate_limit.Policy) (core_rate_limit.Decision, error)
}
