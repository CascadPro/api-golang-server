package outbox_media_handler

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/s3/pool"
	media_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/media/repository/postgres"
)

type Handler struct {
	storage    core_s3_pool.Pool
	repository media_postgres_repository.RepositoryMethods
}

type HandlerMethods interface {
	HandleDeleteFile(ctx context.Context, event domain.OutboxEvent) error
	HandleProcessAvatar(ctx context.Context, event domain.OutboxEvent) error
}

func New(storage core_s3_pool.Pool, repository media_postgres_repository.RepositoryMethods) *Handler {
	return &Handler{
		storage:    storage,
		repository: repository,
	}
}
