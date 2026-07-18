package users_http_dto

import (
	"fmt"
	"strings"

	core_http_types "github.com/Svat-dev/golang-todo/internal/core/transport/http/types"
)

type PatchUserRequestSwagger struct {
	FN string `json:"full_name" example:"John Doe"`
	PH string `json:"phone_number" example:"+79999999999"`
}

type PatchUserRequest struct {
	FullName    core_http_types.Nullable[string] `json:"full_name"`
	PhoneNumber core_http_types.Nullable[string] `json:"phone_number"`
}

func (r *PatchUserRequest) Validate() error {
	if r.FullName.Set {
		if err := validateFullName(r.FullName.Value); err != nil {
			return err
		}
	}

	if r.PhoneNumber.Set {
		if err := validatePhoneNumber(r.PhoneNumber.Value); err != nil {
			return err
		}
	}

	return nil
}

func validateFullName(v *string) error {
	if v == nil {
		return fmt.Errorf("`FullName` can't be NULL")
	}
	nameLen := len([]rune(*v))
	if nameLen < 2 || nameLen > 100 {
		return fmt.Errorf("`FullName` must be between 2 and 100 symbols")
	}
	return nil
}

func validatePhoneNumber(v *string) error {
	if v == nil {
		return nil
	}
	numLen := len([]rune(*v))
	if numLen < 10 || numLen > 15 {
		return fmt.Errorf("`PhoneNumber` must be between 10 and 15 symbols")
	}
	if !strings.HasPrefix(*v, "+") {
		return fmt.Errorf("`PhoneNumber` must starts with `+`")
	}
	return nil
}
