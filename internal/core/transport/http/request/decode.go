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

func DecodeAndValidate(rw http.ResponseWriter, r *http.Request, dest any) error {
	const maxJSONBody = 1 << 20 // 1 MiB

	r.Body = http.MaxBytesReader(rw, r.Body, maxJSONBody)

	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(dest); err != nil {
		return fmt.Errorf("decode JSON: %w: %w", err, core_errors.ErrInvalidArgument)
	}

	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		return fmt.Errorf("unexpected data after JSON body: %w", core_errors.ErrInvalidArgument)
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
