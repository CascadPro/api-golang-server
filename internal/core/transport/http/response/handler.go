package core_http_response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_i18n "github.com/CascadePro/api-golang-server/internal/core/i18n"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

type ResponseHandler struct {
	log    *core_logger.Logger
	locale language.Tag
	rw     http.ResponseWriter
}

func NewResponseHandler(log *core_logger.Logger, locale language.Tag, rw http.ResponseWriter) *ResponseHandler {
	return &ResponseHandler{
		log:    log,
		locale: locale,
		rw:     rw,
	}
}

func (h *ResponseHandler) MediaContentResponse(bytes []byte, mimeType domain.FileMimeType) {
	h.rw.Header().Add("Content-Type", string(mimeType))
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
	descriptor := mapError(err)

	h.logError(descriptor, err, msg)

	descriptor.Message = core_i18n.Translate(
		h.locale,
		string(descriptor.Code),
		"Internal server error",
	)

	h.errorResponse(descriptor)
}

func (h *ResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	descriptor := ErrorDescriptor{
		Code:       core_errors.CodeInternal,
		StatusCode: statusCode,
		Message:    msg,
	}

	h.errorResponse(descriptor)
}

func (h *ResponseHandler) errorResponse(descriptor ErrorDescriptor) {
	body := ErrorResponse{
		Code:      string(descriptor.Code),
		Message:   descriptor.Message,
		Fields:    descriptor.Fields,
		Timestamp: time.Now().UTC(),
	}

	h.JsonResponse(body, descriptor.StatusCode)
}

func (h *ResponseHandler) logError(descriptor ErrorDescriptor, err error, msg string) {
	switch descriptor.StatusCode {
	case http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusConflict,
		http.StatusTooManyRequests:

		h.log.Warn(
			msg,
			zap.String("code", string(descriptor.Code)),
			zap.Error(err),
		)

	case http.StatusNotFound:

		h.log.Debug(
			msg,
			zap.String("code", string(descriptor.Code)),
			zap.Error(err),
		)

	default:

		h.log.Error(
			msg,
			zap.String("code", string(descriptor.Code)),
			zap.Error(err),
		)
	}
}
