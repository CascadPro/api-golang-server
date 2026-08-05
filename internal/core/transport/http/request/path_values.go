package core_http_request

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

func GetUUIDPathValue(r *http.Request, key string) (uuid.UUID, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return uuid.Nil, fmt.Errorf("no key='%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	value, err := uuid.Parse(pathValue)
	if err != nil {
		return uuid.Nil, fmt.Errorf("key='%s' in path values is not valid UUID: %w", key, err)
	}
	if value == uuid.Nil {
		return uuid.Nil, fmt.Errorf("key='%s' in path values can't be NULL: %w", key, core_errors.ErrInvalidArgument)
	}

	return value, nil
}

func GetIDPathValue(r *http.Request, key string, length int) (string, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return "", fmt.Errorf("no key='%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	if err := core_validation.ValidateID(pathValue, length); err != nil {
		return "", fmt.Errorf("key='%s' in path values is not valid ID: %w", key, err)
	}

	return pathValue, nil
}

func GetFileNamePathValue(r *http.Request, key string) (string, *string, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return "", nil, fmt.Errorf("no key='%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	var id, ext string
	filename := strings.SplitN(pathValue, ".", 2)

	id = filename[0]

	if err := core_validation.ValidateID(id, domain.FileIDByteLength); err != nil {
		return "", nil, fmt.Errorf("validate filename: %w", err)
	}

	if len(filename) > 1 {
		ext = filename[1]
	}

	return id, &ext, nil
}

func GetFileTagPathValue(r *http.Request, key string) (domain.FileTag, error) {
	pathValue := r.PathValue(key)
	if pathValue == "" {
		return domain.FileTagNil, fmt.Errorf("no key='%s' in path values: %w", key, core_errors.ErrInvalidArgument)
	}

	fileTag := domain.FileTag(pathValue)

	if err := core_validation.ValidateArray(domain.FileTags, fileTag); err != nil {
		return domain.FileTagNil, fmt.Errorf("key=%s is not valid `FileTag`: %w", key, err)
	}

	return fileTag, nil
}
