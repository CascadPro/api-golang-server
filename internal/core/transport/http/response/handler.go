package core_http_response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"go.uber.org/zap"
)

type ResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewResponseHandler(log *core_logger.Logger, rw http.ResponseWriter) *ResponseHandler {
	return &ResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *ResponseHandler) MediaContentResponse(bytes []byte, contentType string) {
	h.rw.Header().Add("Content-Type", contentType)
	h.rw.Header().Add("Content-Length", fmt.Sprintf("%d", len(bytes)))

	h.rw.WriteHeader(http.StatusOK)

	h.rw.Write(bytes)
}

func (h *ResponseHandler) NoContentResponse() {
	h.rw.WriteHeader(http.StatusNoContent)
}

func (h *ResponseHandler) JsonResponse(body any, statusCode int) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(body); err != nil {
		h.log.Error("write http response", zap.Error(err))
	}
}

func (h *ResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrUnauthorized):
		statusCode = http.StatusUnauthorized
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrForbidden):
		statusCode = http.StatusForbidden
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug

	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn

	case errors.Is(err, core_errors.ErrTooManyRequests):
		statusCode = http.StatusTooManyRequests
		logFunc = h.log.Warn

	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *ResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(statusCode, err, msg)
}

func (h *ResponseHandler) errorResponse(statusCode int, err error, msg string) {
	body := ErrorResponse{
		Error:     err.Error(),
		Message:   msg,
		Timestamp: time.Now().UTC(),
	}

	h.JsonResponse(body, statusCode)
}
