package media_postgres_repository

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
)

type Repository struct {
	pool core_postgres_pool.Pool
}

type RepositoryMethods interface {
	CreateFile(ctx context.Context, file *domain.File) (domain.File, error)
	GetFile(ctx context.Context, fileID string) (domain.File, []byte, error)

	DeleteFile(ctx context.Context, fileID string) error
	DeleteFileTx(ctx context.Context, tx core_postgres_pool.Tx, fileID string) error

	PatchFile(ctx context.Context, fileID string, file domain.File) (domain.File, error)
	PatchFileTx(ctx context.Context, tx core_postgres_pool.Tx, fileID string, file domain.File) (domain.File, error)

	BeginTx(ctx context.Context) (core_postgres_pool.Tx, error)
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func (r *Repository) BeginTx(ctx context.Context) (core_postgres_pool.Tx, error) {
	return r.pool.Begin(ctx)
}
