package core_ws_response

import (
	"encoding/json"
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_i18n "github.com/CascadePro/api-golang-server/internal/core/i18n"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_ws_conn "github.com/CascadePro/api-golang-server/internal/core/transport/ws/conn"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"golang.org/x/text/language"
)

type Handler struct {
	conn   *core_ws_conn.Conn
	logger *core_logger.Logger
	locale language.Tag
}

func NewHandler(conn *core_ws_conn.Conn, logger *core_logger.Logger, locale language.Tag) *Handler {
	return &Handler{
		conn:   conn,
		logger: logger,
		locale: locale,
	}
}

func (h *Handler) Response(body any) {
	h.conn.WriteJSON(body)
}

func (h *Handler) ErrorResponse(err error, msg string) {
	descriptor := mapError(err)

	h.logError(descriptor, err, msg)

	descriptor.Message = core_i18n.Translate(
		h.locale,
		string(descriptor.Code),
		"Internal server error",
	)

	h.errorResponse(descriptor)
}

func (h *Handler) PanicResponse(p any, msg string) {
	statusCode := websocket.CloseInternalServerErr
	err := fmt.Errorf("unexpected panic: %v", p)

	h.logger.Error(msg, zap.Error(err))

	descriptor := ErrorDescriptor{
		Code:       core_errors.CodeInternal,
		StatusCode: statusCode,
		Message:    msg,
	}

	h.errorResponse(descriptor)
}

func (h *Handler) errorResponse(descriptor ErrorDescriptor) {
	body := ErrorResponse{
		Code:      descriptor.Code,
		Message:   descriptor.Message,
		Timestamp: time.Now().Unix(),
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		h.logger.Error("failed to marshal error response", zap.Error(err))
		return
	}

	closePayload := websocket.FormatCloseMessage(descriptor.StatusCode, string(jsonBody))

	h.conn.WriteControl(websocket.CloseMessage, closePayload, time.Now().Add(time.Second))
}

func (h *Handler) logError(descriptor ErrorDescriptor, err error, msg string) {
	switch descriptor.StatusCode {
	case websocket.CloseInvalidFramePayloadData,
		websocket.CloseUnsupportedData,
		websocket.CloseNoStatusReceived,
		websocket.CloseTryAgainLater:

		h.logger.Warn(
			msg,
			zap.String("code", string(descriptor.Code)),
			zap.Error(err),
		)

	default:

		h.logger.Error(
			msg,
			zap.String("code", string(descriptor.Code)),
			zap.Error(err),
		)
	}
}

func (h *Handler) AuthErrorResponse(err error) {
	descriptor := mapError(err)

	message := core_i18n.Translate(h.locale, string(descriptor.Code), "Authentication failed")

	body := NewAuthErrorResponse(descriptor.Code, message)

	if err := h.conn.WriteJSON(body); err != nil {
		h.logger.Error("failed to write websocket auth error", zap.Error(err))
		return
	}

	_ = h.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(descriptor.StatusCode, ""),
		time.Now().Add(time.Second),
	)
}
