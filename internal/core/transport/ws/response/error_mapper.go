package core_ws_response

import (
	"errors"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/gorilla/websocket"
)

func mapError(err error) ErrorDescriptor {
	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		return ErrorDescriptor{
			Code:       core_errors.CodeInvalidArgument,
			StatusCode: websocket.CloseInvalidFramePayloadData,
		}

	case errors.Is(err, core_errors.ErrTooManyRequests):
		return ErrorDescriptor{
			Code:       core_errors.CodeTooManyRequests,
			StatusCode: websocket.CloseTryAgainLater,
		}

	default:
		return ErrorDescriptor{
			Code:       core_errors.CodeInternal,
			StatusCode: websocket.CloseInternalServerErr,
		}
	}
}
