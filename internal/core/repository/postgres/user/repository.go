package core_postgres_user

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
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
