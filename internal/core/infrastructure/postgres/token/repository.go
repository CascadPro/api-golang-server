package core_postgres_token

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	"github.com/google/uuid"
)

type Repository struct {
	pool core_postgres_pool.Pool
}

type RepositoryMethods interface {
	CreateToken(ctx context.Context, token domain.Token) (domain.Token, error)
	GetToken(ctx context.Context, token domain.Token) (domain.Token, error)
	DeleteToken(ctx context.Context, id uuid.UUID) error
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
