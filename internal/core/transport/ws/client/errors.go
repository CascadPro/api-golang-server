package core_ws_client

const (
	ErrorCodeInvalidMessage  = "invalid_message"
	ErrorCodeUnauthorized    = "unauthorized"
	ErrorCodeInvalidToken    = "invalid_token"
	ErrorCodeSessionNotFound = "session_not_found"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
