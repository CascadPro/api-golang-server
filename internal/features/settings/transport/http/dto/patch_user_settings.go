package settings_http_dto

import (
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_http_types "github.com/CascadePro/api-golang-server/internal/core/transport/http/types"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

type PatchUserSettingsRequest struct {
	SessionExpireTerm core_http_types.Nullable[domain.SessionExpireTime] `json:"session_expire_term" example:"30d"`
}

func (r *PatchUserSettingsRequest) Validate() error {
	if r.SessionExpireTerm.Set {
		if r.SessionExpireTerm.Value == nil || *r.SessionExpireTerm.Value == domain.SessionExpireNil {
			return fmt.Errorf("`SessionExpireTerm` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
		}

		if err := core_validation.ValidateArray(domain.SessionExpireTimes, *r.SessionExpireTerm.Value); err != nil {
			return fmt.Errorf("`SessionExpireTerm` isn't valid enum type: %w", err)
		}
	}

	return nil
}
