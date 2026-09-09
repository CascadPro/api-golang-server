package core_http_response

import (
	"errors"
	"net/http"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	auth_errors "github.com/CascadePro/api-golang-server/internal/features/auth/errors"
	client_errors "github.com/CascadePro/api-golang-server/internal/features/client/errors"
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
	requests_errors "github.com/CascadePro/api-golang-server/internal/features/requests/errors"
	sessions_errors "github.com/CascadePro/api-golang-server/internal/features/sessions/errors"
	settings_errors "github.com/CascadePro/api-golang-server/internal/features/settings/errors"
	users_errors "github.com/CascadePro/api-golang-server/internal/features/users/errors"
)

type ErrorDescriptor struct {
	Code       core_errors.Code
	StatusCode int
	Message    string
	Fields     map[string]string
}

func mapError(err error) ErrorDescriptor {
	var validationErr *core_errors.ValidationError

	if errors.As(err, &validationErr) {
		return ErrorDescriptor{
			Code:       core_errors.CodeValidation,
			StatusCode: http.StatusBadRequest,
			Fields:     validationErr.Fields,
		}
	}

	switch {
	case errors.Is(err, auth_errors.ErrInvalidCredentials):
		return ErrorDescriptor{
			Code:       core_errors.CodeInvalidCredentials,
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, auth_errors.ErrInvalidRefreshToken):
		return ErrorDescriptor{
			Code:       core_errors.CodeInvalidToken,
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, auth_errors.ErrSessionExpired):
		return ErrorDescriptor{
			Code:       core_errors.CodeTokenExpired,
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, users_errors.ErrUserNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeUserNotFound,
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, users_errors.ErrUserNotActivated):
		return ErrorDescriptor{
			Code:       core_errors.CodeUserNotActivated,
			StatusCode: http.StatusForbidden,
		}

	case errors.Is(err, users_errors.ErrEmailAlreadyExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeEmailAlreadyExists,
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, users_errors.ErrUsernameExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeUsernameExists,
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, users_errors.ErrAvatarNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileNotFound,
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, client_errors.ErrClientNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeClientNotFound,
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, client_errors.ErrClientAlreadyExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeClientAlreadyExists,
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, media_errors.ErrFileNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileNotFound,
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, media_errors.ErrFileAlreadyExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileAlreadyExists,
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, media_errors.ErrFileAccessDenied):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileAccessDenied,
			StatusCode: http.StatusForbidden,
		}

	case errors.Is(err, media_errors.ErrUnsupportedFile):
		return ErrorDescriptor{
			Code:       core_errors.CodeUnsupportedFile,
			StatusCode: 400,
		}

	case errors.Is(err, media_errors.ErrFileTooLarge):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileTooLarge,
			StatusCode: 413,
		}

	case errors.Is(err, requests_errors.ErrRequestNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeRequestNotFound,
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, requests_errors.ErrRequestAlreadyExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeRequestAlreadyExists,
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, requests_errors.ErrRequestAccessDenied):
		return ErrorDescriptor{
			Code:       core_errors.CodeRequestAccessDenied,
			StatusCode: http.StatusForbidden,
		}

	case errors.Is(err, requests_errors.ErrRequestAlreadyHandled):
		return ErrorDescriptor{
			Code:       core_errors.CodeRequestAlreadyHandled,
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, sessions_errors.ErrSessionNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeSessionNotFound,
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, sessions_errors.ErrSessionExpired):
		return ErrorDescriptor{
			Code:       core_errors.CodeSessionExpired,
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, sessions_errors.ErrSessionRevoked):
		return ErrorDescriptor{
			Code:       core_errors.CodeSessionRevoked,
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, settings_errors.ErrSettingsNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeSettingsNotFound,
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, core_errors.ErrInvalidArgument):
		return ErrorDescriptor{
			Code:       core_errors.CodeInvalidArgument,
			StatusCode: http.StatusBadRequest,
		}

	case errors.Is(err, core_errors.ErrUnauthorized):
		return ErrorDescriptor{
			Code:       core_errors.CodeUnauthorized,
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, core_errors.ErrForbidden):
		return ErrorDescriptor{
			Code:       core_errors.CodeForbidden,
			StatusCode: http.StatusForbidden,
		}

	case errors.Is(err, core_errors.ErrNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeNotFound,
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, core_errors.ErrConflict):
		return ErrorDescriptor{
			Code:       core_errors.CodeConflict,
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, core_errors.ErrTooManyRequests):
		return ErrorDescriptor{
			Code:       core_errors.CodeTooManyRequests,
			StatusCode: http.StatusTooManyRequests,
		}

	default:
		return ErrorDescriptor{
			Code:       core_errors.CodeInternal,
			StatusCode: http.StatusInternalServerError,
		}
	}
}
