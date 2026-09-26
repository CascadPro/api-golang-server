package outbox_realtime_handler

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_redis_realtime "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/realtime"
)

type Handler struct {
	publisher core_redis_realtime.PublisherMethods
}

type HandlerMethods interface {
	HandlePublishEvent(ctx context.Context, event domain.OutboxEvent) error
}

func New(publisher core_redis_realtime.PublisherMethods) *Handler {
	return &Handler{
		publisher: publisher,
	}
}
