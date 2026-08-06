package client_postgres_repository

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
	CreateClient(ctx context.Context, client domain.Client) (domain.Client, error)
	GetClient(ctx context.Context, id uuid.UUID) (domain.Client, error)
	DeleteClient(ctx context.Context, id uuid.UUID) error
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
