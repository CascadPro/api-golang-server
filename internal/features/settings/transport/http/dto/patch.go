package settings_http_dto

import (
	"fmt"
	"slices"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_http_types "github.com/CascadePro/api-golang-server/internal/core/transport/http/types"
)

type PatchUserSettingsRequest struct {
	Version           int64                                              `json:"version"             swaggertype:"number"`
	SessionExpireTerm core_http_types.Nullable[domain.SessionExpireTime] `json:"session_expire_term" swaggertype:"string"`
	Session2FAEnabled core_http_types.Nullable[bool]                     `json:"session_2fa_enabled" swaggertype:"boolean"`
	Preferences       PatchUserSettingsPreferences                       `json:"preferences"`
	Notifications     PatchUserSettingsNotifications                     `json:"notifications"`
}

type PatchUserSettingsPreferences struct {
	Theme                   core_http_types.Nullable[domain.PreferencesTheme]     `json:"theme"                     swaggertype:"string"`
	Accent                  core_http_types.Nullable[domain.PreferencesAccent]    `json:"accent"                    swaggertype:"string"`
	Animations              core_http_types.Nullable[domain.PreferencesAnimation] `json:"animations"                swaggertype:"string"`
	Density                 core_http_types.Nullable[domain.PreferencesDensity]   `json:"density"                   swaggertype:"string"`
	TextSize                core_http_types.Nullable[domain.PreferencesTextSize]  `json:"text_size"                 swaggertype:"string"`
	ConfirmImportantActions core_http_types.Nullable[bool]                        `json:"confirm_important_actions" swaggertype:"boolean"`
}

type PatchUserSettingsNotifications struct {
	Requests core_http_types.Nullable[bool] `json:"requests" swaggertype:"boolean"`
	Sessions core_http_types.Nullable[bool] `json:"sessions" swaggertype:"boolean"`
	System   core_http_types.Nullable[bool] `json:"system"   swaggertype:"boolean"`
	Push     core_http_types.Nullable[bool] `json:"push"     swaggertype:"boolean"`
}

func validateNullable[T any](name string, value core_http_types.Nullable[T]) error {
	if value.Set && value.Value == nil {
		return fmt.Errorf("`%s` can't be patched to NULL: %w", name, core_errors.ErrInvalidArgument)
	}
	return nil
}

func (r *PatchUserSettingsRequest) Validate() error {
	if r.Version < 1 {
		return fmt.Errorf("`Version` must be greater than zero: %w", core_errors.ErrInvalidArgument)
	}
	if !r.HasChanges() {
		return fmt.Errorf("settings patch is empty: %w", core_errors.ErrInvalidArgument)
	}

	if err := validateNullable("SessionExpireTerm", r.SessionExpireTerm); err != nil {
		return err
	}
	if r.SessionExpireTerm.Set && !slices.Contains(domain.SessionExpireTimes, *r.SessionExpireTerm.Value) {
		return fmt.Errorf("invalid `SessionExpireTerm`: %w", core_errors.ErrInvalidArgument)
	}
	if err := validateNullable("Session2FAEnabled", r.Session2FAEnabled); err != nil {
		return err
	}

	p := r.Preferences
	if err := validateNullable("Preferences.Theme", p.Theme); err != nil {
		return err
	}
	if p.Theme.Set && !slices.Contains(domain.PreferencesThemes, *p.Theme.Value) {
		return fmt.Errorf("invalid `Preferences.Theme`: %w", core_errors.ErrInvalidArgument)
	}
	if err := validateNullable("Preferences.Accent", p.Accent); err != nil {
		return err
	}
	if p.Accent.Set && !slices.Contains(domain.PreferencesAccents, *p.Accent.Value) {
		return fmt.Errorf("invalid `Preferences.Accent`: %w", core_errors.ErrInvalidArgument)
	}
	if err := validateNullable("Preferences.Animations", p.Animations); err != nil {
		return err
	}
	if p.Animations.Set && !slices.Contains(domain.PreferencesAnimations, *p.Animations.Value) {
		return fmt.Errorf("invalid `Preferences.Animations`: %w", core_errors.ErrInvalidArgument)
	}
	if err := validateNullable("Preferences.Density", p.Density); err != nil {
		return err
	}
	if p.Density.Set && !slices.Contains(domain.PreferencesDensities, *p.Density.Value) {
		return fmt.Errorf("invalid `Preferences.Density`: %w", core_errors.ErrInvalidArgument)
	}
	if err := validateNullable("Preferences.TextSize", p.TextSize); err != nil {
		return err
	}
	if p.TextSize.Set && !slices.Contains(domain.PreferencesTextSizes, *p.TextSize.Value) {
		return fmt.Errorf("invalid `Preferences.TextSize`: %w", core_errors.ErrInvalidArgument)
	}
	if err := validateNullable("Preferences.ConfirmImportantActions", p.ConfirmImportantActions); err != nil {
		return err
	}

	n := r.Notifications
	for _, item := range []struct {
		name  string
		value core_http_types.Nullable[bool]
	}{
		{"Notifications.Requests", n.Requests},
		{"Notifications.Sessions", n.Sessions},
		{"Notifications.System", n.System},
		{"Notifications.Push", n.Push},
	} {
		if err := validateNullable(item.name, item.value); err != nil {
			return err
		}
	}
	return nil
}

func (r *PatchUserSettingsRequest) HasChanges() bool {
	p, n := r.Preferences, r.Notifications
	return r.SessionExpireTerm.Set || r.Session2FAEnabled.Set ||
		p.Theme.Set || p.Accent.Set || p.Animations.Set || p.Density.Set ||
		p.TextSize.Set || p.ConfirmImportantActions.Set ||
		n.Requests.Set || n.Sessions.Set || n.System.Set || n.Push.Set
}
