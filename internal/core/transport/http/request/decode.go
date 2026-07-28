package core_http_request

import (
	"encoding/json"
	"fmt"
	"net/http"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var err error

	v, ok := dest.(core_validation.Validatable)
	if ok {
		err = v.Validate()
	} else {
		if e := core_validation.Validation.Struct(dest); e != nil {
			errs := core_validation.MapErrors(e)
			err = errs[0]
		}
	}

	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
