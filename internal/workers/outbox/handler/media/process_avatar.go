package outbox_media_handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (h *Handler) HandleProcessAvatar(ctx context.Context, event domain.OutboxEvent) error {
	var payload domain.EventMediaAvatarProcessPayload

	if err := json.Unmarshal(event.Payload, &payload); err != nil {
		return fmt.Errorf("decode media avatar process payload: %w", err)
	}

	if err := core_validation.ValidateID(payload.FileID, domain.FileIDByteLength); err != nil {
		return fmt.Errorf("validate file id: %w", err)
	}

	avatar, _, err := h.repository.GetFile(ctx, payload.FileID)
	if err != nil {
		return fmt.Errorf("get avatar from repository: %w", err)
	}

	key := fmt.Sprintf("%s/%s", domain.FileTagAvatars, payload.FileID)
	content, err := h.storage.GetObject(ctx, key)
	if err != nil {
		return fmt.Errorf("get avatar from storage: %w", err)
	}

	if err := avatar.GeneratePlaceholder(content); err != nil {
		return fmt.Errorf("generate placeholder: %w", err)
	}

	if err := h.repository.PatchPlaceholder(ctx, payload.FileID, avatar.Version, avatar.GetPlaceholder()); err != nil {
		return fmt.Errorf("update avatar in repository: %w", err)
	}

	return nil
}
