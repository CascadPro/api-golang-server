package core_ipinfo_client

import (
	core_ipinfo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/ipinfo/pool"
)

type Repository struct {
	pool core_ipinfo_pool.Pool
}

func NewRepository(pool core_ipinfo_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
