package outbox_media_handler

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/s3/pool"
)

type Handler struct {
	storage core_s3_pool.Pool
}

type HandlerMethods interface {
	HandleDeleteFile(ctx context.Context, event domain.OutboxEvent) error
}

func New(storage core_s3_pool.Pool) *Handler {
	return &Handler{
		storage: storage,
	}
}
