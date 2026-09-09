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
			Message:    "Проверьте правильность введённых данных",
			Fields:     validationErr.Fields,
		}
	}

	switch {
	case errors.Is(err, auth_errors.ErrInvalidCredentials):
		return ErrorDescriptor{
			Code:       core_errors.CodeInvalidCredentials,
			Message:    "Неверная почта или пароль",
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, auth_errors.ErrInvalidRefreshToken):
		return ErrorDescriptor{
			Code:       core_errors.CodeInvalidToken,
			Message:    "Недействительный токен",
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, auth_errors.ErrSessionExpired):
		return ErrorDescriptor{
			Code:       core_errors.CodeTokenExpired,
			Message:    "Срок действия токена истёк",
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, users_errors.ErrUserNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeUserNotFound,
			Message:    "Пользователь не найден",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, users_errors.ErrUserNotActivated):
		return ErrorDescriptor{
			Code:       core_errors.CodeUserNotActivated,
			Message:    "Пользователь не активирован",
			StatusCode: http.StatusForbidden,
		}

	case errors.Is(err, users_errors.ErrEmailAlreadyExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeEmailAlreadyExists,
			Message:    "Пользователь с такой почтой уже существует",
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, users_errors.ErrUsernameExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeUsernameExists,
			Message:    "Это имя пользователя уже занято",
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, users_errors.ErrAvatarNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileNotFound,
			Message:    "Аватар не найден",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, client_errors.ErrClientNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeClientNotFound,
			Message:    "Клиент не найден",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, client_errors.ErrClientAlreadyExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeClientAlreadyExists,
			Message:    "Клиент уже существует",
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, media_errors.ErrFileNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileNotFound,
			Message:    "Файл не найден",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, media_errors.ErrFileAlreadyExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileAlreadyExists,
			Message:    "Файл уже существует",
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, media_errors.ErrFileAccessDenied):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileAccessDenied,
			Message:    "Нет доступа к файлу",
			StatusCode: http.StatusForbidden,
		}

	case errors.Is(err, media_errors.ErrUnsupportedFile):
		return ErrorDescriptor{
			Code:       core_errors.CodeUnsupportedFile,
			Message:    "Неподдерживаемый тип файла",
			StatusCode: 400,
		}

	case errors.Is(err, media_errors.ErrFileTooLarge):
		return ErrorDescriptor{
			Code:       core_errors.CodeFileTooLarge,
			Message:    "Файл слишком большой",
			StatusCode: 413,
		}

	case errors.Is(err, requests_errors.ErrRequestNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeRequestNotFound,
			Message:    "Запрос не найден",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, requests_errors.ErrRequestAlreadyExists):
		return ErrorDescriptor{
			Code:       core_errors.CodeRequestAlreadyExists,
			Message:    "Такой запрос уже существует",
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, requests_errors.ErrRequestAccessDenied):
		return ErrorDescriptor{
			Code:       core_errors.CodeRequestAccessDenied,
			Message:    "Нет доступа к запросу",
			StatusCode: http.StatusForbidden,
		}

	case errors.Is(err, requests_errors.ErrRequestAlreadyHandled):
		return ErrorDescriptor{
			Code:       core_errors.CodeRequestAlreadyHandled,
			Message:    "Запрос уже обработан",
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, sessions_errors.ErrSessionNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeSessionNotFound,
			Message:    "Сессия не найдена",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, sessions_errors.ErrSessionExpired):
		return ErrorDescriptor{
			Code:       core_errors.CodeSessionExpired,
			Message:    "Сессия истекла",
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, sessions_errors.ErrSessionRevoked):
		return ErrorDescriptor{
			Code:       core_errors.CodeSessionRevoked,
			Message:    "Сессия отозвана",
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, settings_errors.ErrSettingsNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeSettingsNotFound,
			Message:    "Настройки не найдены",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, core_errors.ErrInvalidArgument):
		return ErrorDescriptor{
			Code:       core_errors.CodeInvalidArgument,
			Message:    "Некорректные данные",
			StatusCode: http.StatusBadRequest,
		}

	case errors.Is(err, core_errors.ErrUnauthorized):
		return ErrorDescriptor{
			Code:       core_errors.CodeUnauthorized,
			Message:    "Требуется авторизация",
			StatusCode: http.StatusUnauthorized,
		}

	case errors.Is(err, core_errors.ErrForbidden):
		return ErrorDescriptor{
			Code:       core_errors.CodeForbidden,
			Message:    "Недостаточно прав",
			StatusCode: http.StatusForbidden,
		}

	case errors.Is(err, core_errors.ErrNotFound):
		return ErrorDescriptor{
			Code:       core_errors.CodeNotFound,
			Message:    "Ресурс не найден",
			StatusCode: http.StatusNotFound,
		}

	case errors.Is(err, core_errors.ErrConflict):
		return ErrorDescriptor{
			Code:       core_errors.CodeConflict,
			Message:    "Конфликт данных",
			StatusCode: http.StatusConflict,
		}

	case errors.Is(err, core_errors.ErrTooManyRequests):
		return ErrorDescriptor{
			Code:       core_errors.CodeTooManyRequests,
			Message:    "Слишком много запросов",
			StatusCode: http.StatusTooManyRequests,
		}

	default:
		return ErrorDescriptor{
			Code:       core_errors.CodeInternal,
			Message:    "Внутренняя ошибка сервера",
			StatusCode: http.StatusInternalServerError,
		}
	}
}
