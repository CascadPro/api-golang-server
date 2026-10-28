package outbox_media_handler

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/s3/pool"
)

type Handler struct {
	storage    core_s3_pool.Pool
	repository Repository
}

type Repository interface {
	GetFile(ctx context.Context, fileID string) (domain.File, []byte, error)

	DeleteFileTx(ctx context.Context, tx core_postgres_pool.Tx, fileID string) error

	PatchPlaceholder(ctx context.Context, fileID string, version int64, placeholder []byte) error
}

type HandlerMethods interface {
	HandleDeleteFile(ctx context.Context, event domain.OutboxEvent) error

	HandleProcessAvatar(ctx context.Context, event domain.OutboxEvent) error
}

func New(storage core_s3_pool.Pool, repository Repository) *Handler {
	return &Handler{
		storage:    storage,
		repository: repository,
	}
}
