package core_ws_dispatcher

import (
	"encoding/json"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	core_ws_conn "github.com/CascadePro/api-golang-server/internal/core/transport/ws/conn"
	core_ws_response "github.com/CascadePro/api-golang-server/internal/core/transport/ws/response"
	presence_service "github.com/CascadePro/api-golang-server/internal/features/presence/service"
)

type Dispatcher struct {
	Presence presence_service.ServiceMethods

	tokenIssuer core_jwt_security.AccessTokenVerifier
}

func NewDispatcher(
	presence presence_service.ServiceMethods,
	tokenIssuer core_jwt_security.AccessTokenVerifier,
) *Dispatcher {
	return &Dispatcher{
		Presence: presence,

		tokenIssuer: tokenIssuer,
	}
}

func (d *Dispatcher) Handler(conn *core_ws_conn.Conn) error {
	ctx := conn.Context()

	log := core_logger.FromContext(ctx)
	locale := core_context.Locale(ctx)

	responseHandler := core_ws_response.NewHandler(conn, log, locale)

	_, data, err := conn.ReadMessage()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to read websocket message")
		return err
	}

	var message Message

	if err := json.Unmarshal(data, &message); err != nil {
		responseHandler.ErrorResponse(err, "failed to unmarshal websocket message")
		return err
	}

	if err := message.Validate(); err != nil {
		responseHandler.ErrorResponse(err, "invalid websocket message")
		return err
	}

	switch message.Type {
	default:
		return fmt.Errorf("unsupported message type: %w", core_errors.ErrInvalidArgument)
	}

	// return nil
}
