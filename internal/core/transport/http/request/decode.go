package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/go-playground/validator/v10"
)

var validation = validator.New()

type validatable interface {
	Validate() error
}

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var err error

	v, ok := dest.(validatable)
	if ok {
		err = v.Validate()
	} else {
		err = validation.Struct(dest)
	}

	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
