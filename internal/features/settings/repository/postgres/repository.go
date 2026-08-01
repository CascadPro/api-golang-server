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
	CreateUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error)
	GetUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error)
	PatchUserSettings(ctx context.Context, id uuid.UUID, settings domain.UserSettings) (domain.UserSettings, error)
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
