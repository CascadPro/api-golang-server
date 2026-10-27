package domain

import (
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

type SessionExpireTime string

const (
	SessionExpire7Days  = SessionExpireTime("7d")
	SessionExpire1Month = SessionExpireTime("30d")
	SessionExpire3Month = SessionExpireTime("90d")
)

var SessionExpireTimes = []SessionExpireTime{
	SessionExpire7Days,
	SessionExpire1Month,
	SessionExpire3Month,
}

type UserSettings struct {
	ID      uuid.UUID
	Version int64

	SessionExpireTerm SessionExpireTime
	Session2FAEnabled bool

	Preferences   UserPreferences
	Notifications UserNotificationPreferences

	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
}

func NewUserSettings(userID uuid.UUID, expireTime SessionExpireTime) UserSettings {
	return UserSettings{
		ID:                UninitializedUUID,
		Version:           UninitializedVersion,
		UserID:            userID,
		SessionExpireTerm: expireTime,
		Preferences:       DefaultUserPreferences(),
		Notifications:     DefaultUserNotificationPreferences(),
	}
}

func (us *UserSettings) Validate() error {
	if err := core_validation.ValidateArray(SessionExpireTimes, us.SessionExpireTerm); err != nil {
		return fmt.Errorf("`SessionExpireTerm` isn't valid enum type: %w", core_errors.ErrInvalidArgument)
	}

	if err := us.Preferences.Validate(); err != nil {
		return fmt.Errorf("validate user preferences: %w", err)
	}

	if us.CreatedAt.After(us.UpdatedAt) {
		return fmt.Errorf("`UpdatedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type UserSettingsPatch struct {
	Version           int64
	SessionExpireTerm Nullable[SessionExpireTime]
	Session2FAEnabled Nullable[bool]
	Preferences       UserPreferencesPatch
	Notifications     UserNotificationPreferencesPatch
}

func NewUserSettingsPatch(
	version int64,
	expireTerm Nullable[SessionExpireTime],
	session2FAEnabled Nullable[bool],
	preferences UserPreferencesPatch,
	notifications UserNotificationPreferencesPatch,
) UserSettingsPatch {
	return UserSettingsPatch{
		Version:           version,
		SessionExpireTerm: expireTerm,
		Session2FAEnabled: session2FAEnabled,
		Preferences:       preferences,
		Notifications:     notifications,
	}
}

func (p *UserSettingsPatch) Validate() error {
	if p.Version < 1 {
		return fmt.Errorf("`Version` must be greater than zero: %w", core_errors.ErrInvalidArgument)
	}

	if p.SessionExpireTerm.Set {
		if p.SessionExpireTerm.Value == nil || *p.SessionExpireTerm.Value == "" {
			return fmt.Errorf("`SessionExpireTerm` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
		}
		if err := core_validation.ValidateArray(SessionExpireTimes, *p.SessionExpireTerm.Value); err != nil {
			return fmt.Errorf("`SessionExpireTerm` isn't valid enum type: %w", core_errors.ErrInvalidArgument)
		}
	}

	if p.Session2FAEnabled.Set && p.Session2FAEnabled.Value == nil {
		return fmt.Errorf("`Session2FAEnabled` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if err := p.Preferences.Validate(); err != nil {
		return fmt.Errorf("validate preferences patch: %w", err)
	}

	if err := p.Notifications.Validate(); err != nil {
		return fmt.Errorf("validate notifications patch: %w", err)
	}

	if !p.HasChanges() {
		return fmt.Errorf("settings patch is empty: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (p *UserSettingsPatch) HasChanges() bool {
	return p.SessionExpireTerm.Set ||
		p.Session2FAEnabled.Set ||
		len(p.Preferences.Changes()) > 0 ||
		len(p.Notifications.Changes()) > 0
}

func (us *UserSettings) ApplyPatch(patch UserSettingsPatch) (UserSettings, error) {
	if err := patch.Validate(); err != nil {
		return UserSettings{}, fmt.Errorf("validate user settings patch: %w", err)
	}

	patched := *us
	if patch.SessionExpireTerm.Set {
		patched.SessionExpireTerm = *patch.SessionExpireTerm.Value
	}
	if patch.Session2FAEnabled.Set {
		patched.Session2FAEnabled = *patch.Session2FAEnabled.Value
	}
	if patch.Preferences.Theme != nil {
		patched.Preferences.Theme = *patch.Preferences.Theme
	}
	if patch.Preferences.Accent != nil {
		patched.Preferences.Accent = *patch.Preferences.Accent
	}
	if patch.Preferences.Animations != nil {
		patched.Preferences.Animations = *patch.Preferences.Animations
	}
	if patch.Preferences.Density != nil {
		patched.Preferences.Density = *patch.Preferences.Density
	}
	if patch.Preferences.TextSize != nil {
		patched.Preferences.TextSize = *patch.Preferences.TextSize
	}
	if patch.Preferences.ConfirmImportantActions != nil {
		patched.Preferences.ConfirmImportantActions = *patch.Preferences.ConfirmImportantActions
	}
	if patch.Notifications.Requests != nil {
		patched.Notifications.Requests = *patch.Notifications.Requests
	}
	if patch.Notifications.Sessions != nil {
		patched.Notifications.Sessions = *patch.Notifications.Sessions
	}
	if patch.Notifications.System != nil {
		patched.Notifications.System = *patch.Notifications.System
	}
	if patch.Notifications.Push != nil {
		patched.Notifications.Push = *patch.Notifications.Push
	}

	return patched, nil
}
