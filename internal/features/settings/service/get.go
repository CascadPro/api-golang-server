package settings_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (s *Service) GetUserSettings(ctx context.Context, userID uuid.UUID) (domain.UserSettings, error) {
	if userID == uuid.Nil {
		return domain.UserSettings{}, fmt.Errorf("`userID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	settings, err := s.settingsPostgresRepo.GetUserSettings(ctx, domain.UserSettings{UserID: userID})
	if err != nil {
		return domain.UserSettings{}, fmt.Errorf("get user settings: %w", err)
	}

	return settings, nil
}
