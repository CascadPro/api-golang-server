package core_ws_middleware

import (
	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_ws_conn "github.com/CascadePro/api-golang-server/internal/core/transport/ws/conn"
	core_ws_response "github.com/CascadePro/api-golang-server/internal/core/transport/ws/response"
)

func Panic() Middleware {
	return func(next Handler) Handler {
		return func(conn *core_ws_conn.Conn) error {
			ctx := conn.Context()
			log := core_logger.FromContext(ctx)
			locale := core_context.Locale(ctx)
			handler := core_ws_response.NewHandler(conn, log, locale)

			defer func() {
				if p := recover(); p != nil {
					handler.PanicResponse(p, "during handle ws request got unexpected panic")
				}
			}()

			return next(conn)
		}
	}
}
