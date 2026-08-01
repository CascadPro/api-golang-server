package settings_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) PatchUserSettings(
	ctx context.Context,
	userID uuid.UUID,
	patch domain.UserSettingsPatch,
) (domain.UserSettings, error) {
	if userID == uuid.Nil {
		return domain.UserSettings{}, fmt.Errorf("`userID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	settings, err := s.GetUserSettings(ctx, userID)
	if err != nil {
		return domain.UserSettings{}, fmt.Errorf("get user settings: %w", err)
	}

	patched, err := settings.ApplyPatch(patch)
	if err != nil {
		return domain.UserSettings{}, fmt.Errorf("apply patch: %w", err)
	}

	settings, err = s.settingsPostgresRepo.PatchUserSettings(ctx, settings.ID, patched)
	if err != nil {
		return domain.UserSettings{}, fmt.Errorf("patch user settings: %w", err)
	}

	return settings, nil
}
