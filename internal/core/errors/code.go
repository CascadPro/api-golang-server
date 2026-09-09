package core_errors

type Code string

const (
	CodeInvalidArgument Code = "invalid_argument"
	CodeValidation      Code = "validation_error"

	CodeUnauthorized    Code = "unauthorized"
	CodeForbidden       Code = "forbidden"
	CodeNotFound        Code = "not_found"
	CodeConflict        Code = "conflict"
	CodeTooManyRequests Code = "too_many_requests"

	CodeInvalidCredentials Code = "invalid_credentials"
	CodeUserNotFound       Code = "user_not_found"
	CodeUserNotActivated   Code = "user_not_activated"

	CodeEmailAlreadyExists Code = "email_already_exists"
	CodeUsernameExists     Code = "username_already_exists"

	CodeClientNotFound      Code = "client_not_found"
	CodeClientAlreadyExists Code = "client_already_exists"

	CodeFileNotFound      Code = "file_not_found"
	CodeFileAlreadyExists Code = "file_already_exists"
	CodeFileAccessDenied  Code = "file_access_denied"
	CodeUnsupportedFile   Code = "unsupported_file"
	CodeFileTooLarge      Code = "file_too_large"

	CodeRequestNotFound       Code = "request_not_found"
	CodeRequestAlreadyExists  Code = "request_already_exists"
	CodeRequestAccessDenied   Code = "request_access_denied"
	CodeRequestAlreadyHandled Code = "request_already_handled"

	CodeSessionNotFound Code = "session_not_found"
	CodeSessionExpired  Code = "session_expired"
	CodeSessionRevoked  Code = "session_revoked"

	CodeSettingsNotFound Code = "settings_not_found"

	CodeInvalidToken Code = "invalid_token"
	CodeTokenExpired Code = "token_expired"

	CodeInternal Code = "internal_error"
)
