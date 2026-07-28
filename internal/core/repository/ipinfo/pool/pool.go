package core_ipinfo_pool

import (
	"context"
	"net"
	"net/http"
	"time"

	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/repository/redis/cache"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	"github.com/ipinfo/go/v2/ipinfo"
)

type Pool interface {
	GetIPInfo(net.IP) (*ipinfo.Core, error)

	OpTimeout() time.Duration
}

type ConnectionPool struct {
	client    *ipinfo.Client
	redisPool *core_redis_pool.Pool
	timeout   time.Duration
}

func NewConnectionPool(ctx context.Context, cfg Config, pool core_redis_pool.Pool) (*ConnectionPool, error) {
	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     30 * time.Second,
		},
		Timeout: 10 * time.Second,
	}

	repo := core_redis_cache.NewRepository(pool, "ipinfo", cfg.CacheTTL)

	cache := ipinfo.NewCache(NewCacheEngine(ctx, repo))

	ipc := ipinfo.NewClient(&httpClient, cache, cfg.Token)

	ipc.BaseURL = &cfg.BaseURL

	return &ConnectionPool{
		client:    ipc,
		redisPool: &pool,
		timeout:   cfg.Timeout,
	}, nil
}
