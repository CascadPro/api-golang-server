package media_postgres_repository

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"
)

type Repository struct {
	pool core_postgres_pool.Pool
}

type RepositoryMethods interface {
	CreateFile(ctx context.Context, file domain.File) (domain.File, error)
	GetFile(ctx context.Context, fileID string) (domain.File, error)
	PatchFile(ctx context.Context, fileID string, file domain.File) (domain.File, error)
	DeleteFile(ctx context.Context, fileID string) error
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
