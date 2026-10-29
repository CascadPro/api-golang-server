package settings_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"github.com/google/uuid"
)

type Service struct {
	settingsPostgresRepo SettingsRepository
	realtimePublisher    RealtimePublisher

	logger *core_logger.Logger
}

type SettingsRepository interface {
	GetUserSettings(ctx context.Context, settings domain.UserSettings) (domain.UserSettings, error)
	PatchUserSettings(ctx context.Context, id uuid.UUID, patch domain.UserSettingsPatch) (domain.UserSettings, error)
}

type RealtimePublisher interface {
	Publish(ctx context.Context, event domain.RealtimeEvent) error
}

func NewService(
	settingsPostgresRepo SettingsRepository,
	realtimePublisher RealtimePublisher,
	logger *core_logger.Logger,
) *Service {
	return &Service{
		settingsPostgresRepo: settingsPostgresRepo,
		realtimePublisher:    realtimePublisher,

		logger: logger,
	}
}
