package request_http_dto

import (
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_http_types "github.com/CascadePro/api-golang-server/internal/core/transport/http/types"
	"github.com/google/uuid"
)

type PatchRequestRequest struct {
	Title     core_http_types.Nullable[string]                       `json:"title"`
	Origin    core_http_types.Nullable[[]PatchRequestRequestOrigin]  `json:"origin"`
	ClientID  core_http_types.Nullable[uuid.UUID]                    `json:"client_id"`
	WorkTypes core_http_types.Nullable[[]string]                     `json:"work_types"`
	Geography core_http_types.Nullable[[]string]                     `json:"geography"`
	Deadline  core_http_types.Nullable[time.Time]                    `json:"deadline"`
}

type PatchRequestRequestSwagger struct {
	Title     string                       `json:"title"`
	Origin    []PatchRequestRequestOrigin  `json:"origin"`
	ClientID  uuid.UUID                    `json:"client_id"`
	WorkTypes []string                     `json:"work_types"`
	Geography []string                     `json:"geography"`
	Deadline  time.Time                    `json:"deadline"`
}

type PatchRequestRequestOrigin struct {
	Type  core_http_types.Nullable[domain.RequestOriginType] `json:"type"`
	Email core_http_types.Nullable[string]                   `json:"email"`
	Phone core_http_types.Nullable[string]                   `json:"phone"`
	Other core_http_types.Nullable[string]                   `json:"other"`
}

func (r *PatchRequestRequest) Validate() error {
	if r.Title.Set && (r.Title.Value == nil || *r.Title.Value == "") {
		return fmt.Errorf("`Title` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if err := r.ValidateOrigin(); err != nil {
		return err
	}

	return nil
}

func (r *PatchRequestRequest) ValidateOrigin() error {
	if !r.Origin.Set {
		return nil
	}

	for i, origin := range *r.Origin.Value {
		if err := validatePatchRequestOrigin(i, origin); err != nil {
			return err
		}
	}

	return nil
}

func validatePatchRequestOrigin(index int, origin PatchRequestRequestOrigin) error {
	defaultErr := "validate `Origin` item at index %d: "

	if origin.Type.Set && origin.Type.Value != nil {
		if err := domain.ValidateRequestOriginType(*origin.Type.Value); err != nil {
			return fmt.Errorf(defaultErr+"`Type` isn't valid enum type: %w", index, core_errors.ErrInvalidArgument)
		}
	}

	// Count how many value fields are set
	valueFieldsSet := 0
	if origin.Email.Set {
		valueFieldsSet++
	}
	if origin.Phone.Set {
		valueFieldsSet++
	}
	if origin.Other.Set {
		valueFieldsSet++
	}

	if valueFieldsSet > 1 {
		return fmt.Errorf(defaultErr+"origin patch can only set one of email, phone, or other: %w", index, core_errors.ErrInvalidArgument)
	}

	if origin.Email.Set && origin.Email.Value != nil {
		if _, err := domain.ParseEmail(*origin.Email.Value); err != nil {
			return fmt.Errorf(defaultErr+"`Email` isn't valid: %w", index, core_errors.ErrInvalidArgument)
		}
	}
	if origin.Phone.Set && origin.Phone.Value != nil {
		if _, err := domain.ParsePhoneNumber(*origin.Phone.Value); err != nil {
			return fmt.Errorf(defaultErr+"`Phone` isn't valid: %w", index, core_errors.ErrInvalidArgument)
		}
	}
	if origin.Other.Set && origin.Other.Value != nil && *origin.Other.Value == "" {
		return fmt.Errorf(defaultErr+"`Other` value cannot be empty: %w", index, core_errors.ErrInvalidArgument)
	}

	// If type is set, ensure corresponding value is set (unless explicitly null to clear)
	if origin.Type.Set && origin.Type.Value != nil {
		switch *origin.Type.Value {
		case domain.RequestOriginTypeEmail:
			if !origin.Email.Set {
				return fmt.Errorf(defaultErr+"`Email` must be set for email type: %w", index, core_errors.ErrInvalidArgument)
			}
		case domain.RequestOriginTypePhone:
			if !origin.Phone.Set {
				return fmt.Errorf(defaultErr+"`Phone` must be set for phone type: %w", index, core_errors.ErrInvalidArgument)
			}
		case domain.RequestOriginTypeOther:
			if !origin.Other.Set {
				return fmt.Errorf(defaultErr+"`Other` must be set for other type: %w", index, core_errors.ErrInvalidArgument)
			}
		}
	}

	return nil
}