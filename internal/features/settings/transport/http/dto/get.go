package settings_http_dto

import (
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type UserSettingsResponse struct {
	ID                uuid.UUID                         `json:"id" example:"00000000-0000-0000-0000-000000000000"`
	Version           int64                             `json:"version" example:"1"`
	SessionExpireTerm domain.SessionExpireTime          `json:"session_expire_term" example:"30d"`
	Session2FAEnabled bool                              `json:"session_2fa_enabled" example:"false"`
	Preferences       UserSettingsResponsePreferences   `json:"preferences"`
	Notifications     UserSettingsResponseNotifications `json:"notifications"`
}

type UserSettingsResponsePreferences struct {
	Theme                   domain.PreferencesTheme     `json:"theme"`
	Accent                  domain.PreferencesAccent    `json:"accent"`
	Animations              domain.PreferencesAnimation `json:"animations"`
	Density                 domain.PreferencesDensity   `json:"density"`
	TextSize                domain.PreferencesTextSize  `json:"text_size"`
	ConfirmImportantActions bool                        `json:"confirm_important_actions"`
}

type UserSettingsResponseNotifications struct {
	Requests bool `json:"requests"`
	Sessions bool `json:"sessions"`
	System   bool `json:"system"`
	Push     bool `json:"push"`
}
