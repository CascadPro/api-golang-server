package core_ipinfo_pool

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"sync"
	"time"

	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/cache"
	"github.com/ipinfo/go/v2/ipinfo"
)

type Pool interface {
	GetIPInfo(net.IP) (*ipinfo.Core, error)

	OpTimeout() time.Duration
}

type ipInfoClient interface {
	GetIPInfo(net.IP) (*ipinfo.Core, error)
}

type ConnectionPool struct {
	client  ipInfoClient
	timeout time.Duration

	flightsMu sync.Mutex
	flights   map[string]*ipInfoCall
}

func New(cfg Config, cacheRepo core_redis_cache.Cache) (*ConnectionPool, error) {
	if cacheRepo == nil {
		return nil, fmt.Errorf("initialize IPInfo: cache dependency is nil")
	}

	httpClient := http.Client{
		Transport: &http.Transport{
			MaxIdleConns:        100,
			MaxIdleConnsPerHost: 20,
			IdleConnTimeout:     30 * time.Second,
		},
		Timeout: cfg.Timeout,
	}

	baseURL, err := url.Parse(cfg.BaseURL)
	if err != nil || (baseURL.Scheme != "http" && baseURL.Scheme != "https") || baseURL.Host == "" {
		return nil, fmt.Errorf("initialize IPInfo: BASEURL must be an absolute HTTP(S) URL")
	}

	ipInfoCache := ipinfo.NewCache(NewCacheEngine(cacheRepo))
	client := ipinfo.NewClient(&httpClient, ipInfoCache, cfg.Token)
	client.BaseURL = baseURL

	return &ConnectionPool{
		client:  client,
		timeout: cfg.Timeout,
		flights: make(map[string]*ipInfoCall),
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.timeout
}
