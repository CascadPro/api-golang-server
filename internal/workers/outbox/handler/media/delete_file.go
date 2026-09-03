package outbox_media_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (h *Handler) HandleDeleteFile(ctx context.Context, event domain.OutboxEvent) error {
	var payload domain.EventTypeMediaDeleteFilePayload

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("decode media delete payload: %w", err)
	}

	if err := core_validation.ValidateArray(domain.FileTags, payload.Tag); err != nil {
		return fmt.Errorf("validate file tag: %w", err)
	}

	if err := core_validation.ValidateID(payload.FileID, domain.FileIDByteLength); err != nil {
		return fmt.Errorf("validate file id: %w", err)
	}

	key := fmt.Sprintf("%s/%s", payload.Tag, payload.FileID)

	if err := h.storage.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("delete object from S3: %w", err)
	}

	return nil
}
