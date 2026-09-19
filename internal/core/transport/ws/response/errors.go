package core_ws_response

import (
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type ErrorResponse struct {
	Code      core_errors.Code `json:"code"`
	Message   string           `json:"message"`
	Timestamp int64            `json:"timestamp"`
}

type ErrorDescriptor struct {
	Code       core_errors.Code
	StatusCode int
	Message    string
}
