package outbox_realtime_handler

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

type Handler struct {
	publisher RealtimePublisher
}

type RealtimePublisher interface {
	Publish(ctx context.Context, event domain.RealtimeEvent) error
}

type HandlerMethods interface {
	HandlePublishEvent(ctx context.Context, event domain.OutboxEvent) error
}

func New(publisher RealtimePublisher) *Handler {
	return &Handler{
		publisher: publisher,
	}
}
