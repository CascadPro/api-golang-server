package worker_outbox

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	outbox_media_handler "github.com/CascadePro/api-golang-server/internal/workers/outbox/handler/media"
)

type Handler struct {
	mediaHandler outbox_media_handler.HandlerMethods
}

type HandlerMethods interface {
	Handle(ctx context.Context, event domain.OutboxEvent) error
}

func NewHandler(mediaHandler outbox_media_handler.HandlerMethods) *Handler {
	return &Handler{
		mediaHandler: mediaHandler,
	}
}

func (h *Handler) Handle(ctx context.Context, event domain.OutboxEvent) error {
	switch event.Type {
	case domain.EventTypeMediaDeleteFile:
		return h.mediaHandler.HandleDeleteFile(ctx, event)

	default:
		return fmt.Errorf("unsupported event type: %s", event.Type)
	}
}
