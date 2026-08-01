package settings_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	settings_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/settings/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	settingsPostgresRepo settings_postgres_repository.RepositoryMethods
}

type ServiceMethods interface {
	GetUserSettings(context.Context, uuid.UUID) (domain.UserSettings, error)
	PatchUserSettings(context.Context, uuid.UUID, domain.UserSettingsPatch) (domain.UserSettings, error)
}

func NewService(settingsPostgresRepo settings_postgres_repository.RepositoryMethods) *Service {
	return &Service{
		settingsPostgresRepo: settingsPostgresRepo,
	}
}
