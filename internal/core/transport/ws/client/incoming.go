package core_ws_client

import (
	"encoding/json"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_ws_conn "github.com/CascadePro/api-golang-server/internal/core/transport/ws/conn"
	core_ws_response "github.com/CascadePro/api-golang-server/internal/core/transport/ws/response"
)

func incomingHandler(conn *core_ws_conn.Conn) error {
	ctx := conn.Context()
	log := core_logger.FromContext(ctx)
	locale := core_context.Locale(ctx)
	responseHandler := core_ws_response.NewHandler(conn, log, locale)

	_, data, err := conn.ReadMessage()
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to read message")
		return err
	}

	var message ClientMessage
	if err := json.Unmarshal(data, &message); err != nil {
		responseHandler.ErrorResponse(err, "failed to unmarshal client message")
		return err
	}

	return nil
}
