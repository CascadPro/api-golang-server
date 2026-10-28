package settings_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type Service struct {
	settingsPostgresRepo SettingsRepository
}

type SettingsRepository interface {
	GetUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error)
	PatchUserSettings(ctx context.Context, id uuid.UUID, patch domain.UserSettingsPatch) (domain.UserSettings, error)
}

type ServiceMethods interface {
	GetUserSettings(context.Context, uuid.UUID) (domain.UserSettings, error)
	PatchUserSettings(context.Context, uuid.UUID, domain.UserSettingsPatch) (domain.UserSettings, error)
}

func NewService(settingsPostgresRepo SettingsRepository) *Service {
	return &Service{
		settingsPostgresRepo: settingsPostgresRepo,
	}
}
