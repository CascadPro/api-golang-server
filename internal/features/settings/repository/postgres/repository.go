package settings_postgres_repository

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"
	"github.com/google/uuid"
)

type Repository struct {
	pool core_postgres_pool.Pool
}

type RepositoryMethods interface {
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
