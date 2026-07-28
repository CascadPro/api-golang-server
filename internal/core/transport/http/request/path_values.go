package core_http_request

import (
	"fmt"
	"net/http"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func GetIDPathValue(r *http.Request, key string, length int) (string, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return "", fmt.Errorf("no key='%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	if err := core_validation.ValidateID(pathValue, length); err != nil {
		return "", fmt.Errorf("key='%s' in path values is not valid ID: %w", key, core_errors.ErrInvalidArgument)
	}

	return pathValue, nil
}
