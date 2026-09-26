package outbox_realtime_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func (h *Handler) HandlePublishEvent(ctx context.Context, event domain.OutboxEvent) error {
	var payload domain.EventRealtimePublishPayload

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("decode realtime publish payload: %w", err)
	}

	realtimeEvent := domain.NewRealtimeEvent(
		payload.Type,
		payload.UserID,
		payload.SessionID,
		payload.Data,
	)

	if err := h.publisher.Publish(ctx, realtimeEvent); err != nil {
		return fmt.Errorf("publish realtime event: %w", err)
	}

	return nil
}
