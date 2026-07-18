package tasks_http_dto

import (
	"fmt"

	core_http_types "github.com/Svat-dev/golang-todo/internal/core/transport/http/types"
)

type PatchTaskRequestSwagger struct {
	T string `json:"title" example:"Another new task"`
	D string `json:"description" example:"Another simple description"`
	C bool   `json:"completed" example:"true"`
}

type PatchTaskRequest struct {
	Title       core_http_types.Nullable[string] `json:"title"`
	Description core_http_types.Nullable[string] `json:"description"`
	Completed   core_http_types.Nullable[bool]   `json:"completed"`
}

func (r *PatchTaskRequest) Validate() error {
	if r.Title.Set {
		if err := validateTitle(r.Title.Value); err != nil {
			return err
		}
	}

	if r.Description.Set {
		if err := validateDescription(r.Description.Value); err != nil {
			return err
		}
	}

	if r.Completed.Set {
		if r.Completed.Value == nil {
			return fmt.Errorf("`Completed` can't be NULL")
		}
	}

	return nil
}

func validateTitle(v *string) error {
	if v == nil {
		return fmt.Errorf("`Title` can't be NULL")
	}
	length := len([]rune(*v))
	if length < 1 || length > 100 {
		return fmt.Errorf("`Title` must be between 1 and 100 symbols")
	}
	return nil
}

func validateDescription(v *string) error {
	if v == nil {
		return nil
	}
	length := len([]rune(*v))
	if length < 1 || length > 1000 {
		return fmt.Errorf("`Description` must be between 1 and 1000 symbols")
	}
	return nil
}
