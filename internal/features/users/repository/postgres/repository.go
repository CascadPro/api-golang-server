package users_postgres_repository

import (
	"context"

	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
)

type Repository struct {
	pool core_postgres_pool.Pool
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) BeginTx(ctx context.Context) (core_postgres_pool.Tx, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	return r.pool.Begin(ctx)
}
