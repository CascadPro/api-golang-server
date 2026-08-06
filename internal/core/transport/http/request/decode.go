package core_http_request

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	core_validation_init "github.com/CascadePro/api-golang-server/internal/core/validation/init"
)

func DecodeAndValidate(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(dest); err != nil {
		if errors.Is(err, io.EOF) {
			return fmt.Errorf("decode json: %v: %w", core_errors.ErrInvalidArgument, core_errors.ErrEmptyRequestBody)
		}

		return fmt.Errorf("decode json: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	var err error

	v, ok := dest.(core_validation.Validatable)
	if ok {
		err = v.Validate()
	} else {
		if e := core_validation_init.Validation.Struct(dest); e != nil {
			errs := core_validation.MapErrors(e)
			err = errs[0]
		}
	}

	if err != nil {
		return fmt.Errorf("request validation: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
