package app

import (
	"fmt"
	"time"

	core_rate_limit "github.com/CascadePro/api-golang-server/internal/core/rate_limit"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_ws_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/ws/middleware"
)

type rateLimitPolicies struct {
	Root          core_rate_limit.Policy
	Media         core_rate_limit.Policy
	Create        core_rate_limit.Policy
	Avatar        core_rate_limit.Policy
	Auth          core_rate_limit.Policy
	WebSocketAuth core_rate_limit.Policy
}

func defaultRateLimitPolicies() rateLimitPolicies {
	return rateLimitPolicies{
		Root:          core_rate_limit.Policy{Name: "http.root", Limit: 100, Window: 5 * time.Minute},
		Media:         core_rate_limit.Policy{Name: "http.media", Limit: 25, Window: 10 * time.Minute},
		Create:        core_rate_limit.Policy{Name: "http.create", Limit: 15, Window: 15 * time.Minute},
		Avatar:        core_rate_limit.Policy{Name: "http.avatar", Limit: 10, Window: 5 * time.Minute},
		Auth:          core_rate_limit.Policy{Name: "http.auth", Limit: 5, Window: time.Minute},
		WebSocketAuth: core_rate_limit.Policy{Name: "ws.auth", Limit: 10, Window: time.Minute},
	}
}

func (a *App) rateLimitOptions() core_rate_limit.Options {
	mode := core_rate_limit.FailureMode(a.cfg.RateLimitFailureMode)
	return core_rate_limit.Options{FailureMode: mode, Observer: a.rateLimitObserver}
}

func (a *App) newHTTPRateLimiter(policy core_rate_limit.Policy) (*core_http_middleware.HTTPRateLimiter, error) {
	options := a.rateLimitOptions()
	return core_http_middleware.NewHTTPRateLimiterWithOptions(a.rateLimiter, policy, core_http_middleware.RateLimitOptions{
		FailureMode: options.FailureMode,
		Observer:    options.Observer,
	})
}

func (a *App) newWSRateLimiter(policy core_rate_limit.Policy) (*core_ws_middleware.WSRateLimiter, error) {
	return core_ws_middleware.NewWSRateLimiterWithOptions(a.rateLimiter, policy, a.rateLimitOptions())
}

func validateRateLimitPolicies(p rateLimitPolicies) error {
	policies := []core_rate_limit.Policy{p.Root, p.Media, p.Create, p.Avatar, p.Auth, p.WebSocketAuth}
	for _, policy := range policies {
		if err := policy.Validate(); err != nil {
			return fmt.Errorf("validate rate limit policy %q: %w", policy.Name, err)
		}
	}
	return nil
}
