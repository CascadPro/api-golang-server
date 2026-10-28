package sessions_redis_repository

import (
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
)

type Repository struct {
	pool core_redis_pool.Pool
}

func NewRepository(pool core_redis_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
