package settings_postgres_repository

import (
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type UserSettingsModel struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Version int

	SessionExpireTerm domain.SessionExpireTime
}

func modelToDomain(model UserSettingsModel) domain.UserSettings {
	return domain.UserSettings{
		ID:                model.ID,
		Version:           model.Version,
		SessionExpireTerm: model.SessionExpireTerm,
		UserID:            model.UserID,
	}
}
