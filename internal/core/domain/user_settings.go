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

var (
	SessionExpireTimes = []SessionExpireTime{SessionExpire7Days, SessionExpire1Month, SessionExpire3Month}
)

type UserSettings struct {
	ID      uuid.UUID
	Version int

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
