package core_ws_response

import (
	"encoding/json"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type Response struct {
	ID        string                   `json:"id"`
	Type      domain.RealtimeEventType `json:"type"`
	Data      json.RawMessage          `json:"data,omitempty"`
	Timestamp time.Time                `json:"timestamp"`
}

func NewResponse(id string, t domain.RealtimeEventType, data json.RawMessage) *Response {
	return &Response{
		ID:        id,
		Type:      t,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
}

type AuthErrorResponse struct {
	Type      domain.RealtimeEventType `json:"type"`
	Code      core_errors.Code         `json:"code"`
	Message   string                   `json:"message"`
	Timestamp time.Time                `json:"timestamp"`
}

func NewAuthErrorResponse(code core_errors.Code, message string) *AuthErrorResponse {
	return &AuthErrorResponse{
		Type:      domain.RealtimeEventAuthError,
		Code:      code,
		Message:   message,
		Timestamp: time.Now().UTC(),
	}
}
