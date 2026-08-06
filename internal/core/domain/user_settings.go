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
	SessionExpireNil    = SessionExpireTime("")
)

var (
	SessionExpireTimes = []SessionExpireTime{SessionExpire7Days, SessionExpire1Month, SessionExpire3Month}
)

type UserSettings struct {
	ID      uuid.UUID
	Version int64

	SessionExpireTerm SessionExpireTime

	CreatedAt time.Time
	UpdatedAt time.Time

	UserID uuid.UUID
}

func NewUserSettings(userID uuid.UUID, expireTime SessionExpireTime) UserSettings {
	return UserSettings{
		ID:                UninitializedUUID,
		Version:           UninitializedVersion,
		UserID:            userID,
		SessionExpireTerm: expireTime,
	}
}

func (us *UserSettings) Validate() error {
	if err := core_validation.ValidateArray(SessionExpireTimes, us.SessionExpireTerm); err != nil {
		return fmt.Errorf("`SessionExpireTerm` isn't valid enum type: %w", core_errors.ErrInvalidArgument)
	}

	if us.CreatedAt.After(us.UpdatedAt) {
		return fmt.Errorf("`UpdatedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type UserSettingsPatch struct {
	SessionExpireTerm Nullable[SessionExpireTime]
}

func NewUserSettingsPatch(expireTerm Nullable[SessionExpireTime]) UserSettingsPatch {
	return UserSettingsPatch{
		SessionExpireTerm: expireTerm,
	}
}

func (p *UserSettingsPatch) Validate() error {
	if p.SessionExpireTerm.Set {
		if p.SessionExpireTerm.Value == nil {
			return fmt.Errorf("`SessionExpireTerm` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
		}
		if err := core_validation.ValidateArray(SessionExpireTimes, *p.SessionExpireTerm.Value); err != nil {
			return fmt.Errorf("`SessionExpireTerm` isn't valid enum type: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}

func (us *UserSettings) ApplyPatch(patch UserSettingsPatch) (UserSettings, error) {
	if err := patch.Validate(); err != nil {
		return UserSettings{}, fmt.Errorf("validate user settings patch: %w", err)
	}

	var patched UserSettings
	patched.Version = us.Version

	if patch.SessionExpireTerm.Set {
		patched.SessionExpireTerm = *patch.SessionExpireTerm.Value
	}

	return patched, nil
}
