package core_ipinfo_client

import (
	"context"
	"net"

	core_ipinfo_pool "github.com/CascadePro/api-golang-server/internal/core/repository/ipinfo/pool"
)

type Repository struct {
	pool core_ipinfo_pool.Pool
}

type RepositoryMethods interface {
}

func NewRepository(pool core_ipinfo_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
