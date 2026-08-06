package client_postgres_repository

import core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"

type Repository struct {
	pool core_postgres_pool.Pool
}

type RepositoryMethods interface {
	// ...
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
