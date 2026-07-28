package core_redis_cache

import (
	"context"
	"time"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
)

type Repository struct {
	pool       core_redis_pool.Pool
	prefix     string
	expiration time.Duration
}

func NewRepository(pool core_redis_pool.Pool, prefix string, expiration time.Duration) *Repository {
	return &Repository{
		pool:       pool,
		prefix:     prefix,
		expiration: expiration,
	}
}
