package settings_postgres_repository

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type UserSettingsModel struct {
	ID      uuid.UUID
	UserID  uuid.UUID
	Version int64

	SessionExpireTerm domain.SessionExpireTime
	Session2FAEnabled bool

	Preferences   json.RawMessage
	Notifications json.RawMessage

	CreatedAt time.Time
	UpdatedAt time.Time
}

type UserPreferencesModel struct {
	Theme                   domain.PreferencesTheme     `json:"0"`
	Accent                  domain.PreferencesAccent    `json:"1"`
	Animations              domain.PreferencesAnimation `json:"2"`
	Density                 domain.PreferencesDensity   `json:"3"`
	TextSize                domain.PreferencesTextSize  `json:"4"`
	ConfirmImportantActions bool                        `json:"5"`
}

func preferencesModelToDomain(model UserPreferencesModel) domain.UserPreferences {
	return domain.UserPreferences{
		Theme:                   model.Theme,
		Accent:                  model.Accent,
		Animations:              model.Animations,
		Density:                 model.Density,
		TextSize:                model.TextSize,
		ConfirmImportantActions: model.ConfirmImportantActions,
	}
}

type UserNotificationPreferencesModel struct {
	Requests bool `json:"0"`
	Sessions bool `json:"1"`
	System   bool `json:"2"`
	Push     bool `json:"3"`
}

func notificationModelToDomain(model UserNotificationPreferencesModel) domain.UserNotificationPreferences {
	return domain.UserNotificationPreferences{
		Requests: model.Requests,
		Sessions: model.Sessions,
		System:   model.System,
		Push:     model.Push,
	}
}

func modelToDomain(model UserSettingsModel) (domain.UserSettings, error) {
	preferences := domain.DefaultUserPreferences()
	preferencesModel := UserPreferencesModel{
		Theme: preferences.Theme, Accent: preferences.Accent,
		Animations: preferences.Animations, Density: preferences.Density,
		TextSize: preferences.TextSize, ConfirmImportantActions: preferences.ConfirmImportantActions,
	}
	if len(model.Preferences) > 0 {
		if err := json.Unmarshal(model.Preferences, &preferencesModel); err != nil {
			return domain.UserSettings{}, fmt.Errorf("unmarshal user preferences: %w", err)
		}
	}

	notifications := domain.DefaultUserNotificationPreferences()
	notificationsModel := UserNotificationPreferencesModel{
		Requests: notifications.Requests, Sessions: notifications.Sessions,
		System: notifications.System, Push: notifications.Push,
	}
	if len(model.Notifications) > 0 {
		if err := json.Unmarshal(model.Notifications, &notificationsModel); err != nil {
			return domain.UserSettings{}, fmt.Errorf("unmarshal user notifications: %w", err)
		}
	}

	return domain.UserSettings{
		ID: model.ID, UserID: model.UserID, Version: model.Version,
		SessionExpireTerm: model.SessionExpireTerm, Session2FAEnabled: model.Session2FAEnabled,
		Preferences:   preferencesModelToDomain(preferencesModel),
		Notifications: notificationModelToDomain(notificationsModel),
		CreatedAt:     model.CreatedAt, UpdatedAt: model.UpdatedAt,
	}, nil
}
