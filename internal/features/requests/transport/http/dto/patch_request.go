package request_http_dto

import (
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_http_types "github.com/CascadePro/api-golang-server/internal/core/transport/http/types"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

type PatchRequestRequest struct {
	Title     core_http_types.Nullable[string]                      `json:"title"`
	Origin    core_http_types.Nullable[[]PatchRequestRequestOrigin] `json:"origin"`
	ClientID  core_http_types.Nullable[uuid.UUID]                   `json:"client_id"`
	WorkTypes core_http_types.Nullable[[]string]                    `json:"work_types"`
	Geography core_http_types.Nullable[[]string]                    `json:"geography"`
	Deadline  core_http_types.Nullable[time.Time]                   `json:"deadline"`
}

type PatchRequestRequestOrigin struct {
	Type  domain.RequestOriginType `json:"type"`
	Value string                   `json:"value"`
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

	if err := core_validation.ValidateArray(domain.RequestOriginTypes, origin.Type); err != nil {
		return fmt.Errorf(defaultErr+"`Type` isn't valid enum type: %w", index, core_errors.ErrInvalidArgument)
	}

	switch origin.Type {
	case domain.RequestOriginTypeEmail:
		if _, err := core_validation.ValidateStringEmail(&origin.Value); err != nil {
			return fmt.Errorf(defaultErr+"`Value` isn't valid email: %w", index, core_errors.ErrInvalidArgument)
		}
	case domain.RequestOriginTypePhone:
		if _, err := core_validation.ValidatePhoneNumber(origin.Value); err != nil {
			return fmt.Errorf(defaultErr+"`Value` isn't valid phone number: %w", index, core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
