package settings_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

func (s *Service) PatchUserSettings(
	ctx context.Context,
	userID uuid.UUID,
	patch domain.UserSettingsPatch,
) (domain.UserSettings, error) {
	if userID == uuid.Nil {
		return domain.UserSettings{}, fmt.Errorf("`userID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if err := patch.Validate(); err != nil {
		return domain.UserSettings{}, fmt.Errorf("validate user settings patch: %w", err)
	}

	current, err := s.GetUserSettings(ctx, userID)
	if err != nil {
		return domain.UserSettings{}, fmt.Errorf("get user settings: %w", err)
	}

	if current.Version != patch.Version {
		return domain.UserSettings{}, fmt.Errorf("user settings version is stale: %w", core_errors.ErrConflict)
	}

	if _, err := current.ApplyPatch(patch); err != nil {
		return domain.UserSettings{}, fmt.Errorf("apply user settings patch: %w", err)
	}

	updated, err := s.settingsPostgresRepo.PatchUserSettings(ctx, current.ID, patch)
	if err != nil {
		return domain.UserSettings{}, fmt.Errorf("patch user settings: %w", err)
	}

	changes := changedSettingsFields(current, updated)

	event := domain.NewRealtimeEvent(domain.RealtimeEventSettingsUpdated, &userID, nil, changes)

	if err := s.realtimePublisher.Publish(ctx, event); err != nil {
		s.logger.Error(
			"failed to publish settings updated event",
			zap.Error(err),
			zap.String("user_id", userID.String()),
		)
	}

	return updated, nil
}
