package core_ws_server

import (
	"context"
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_ws_client "github.com/CascadePro/api-golang-server/internal/core/transport/ws/client"
	core_ws_conn "github.com/CascadePro/api-golang-server/internal/core/transport/ws/conn"
	core_ws_response "github.com/CascadePro/api-golang-server/internal/core/transport/ws/response"
)

func (s *Server) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	logger := core_logger.FromContext(ctx)
	locale := core_context.Locale(ctx)

	responseHandler := core_http_response.NewResponseHandler(logger, locale, rw)

	ip, err := core_context.ClientIP(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to extract client ip")
		return
	}

	conn, err := s.upgrader.Upgrade(rw, r, nil)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to upgrade connection")
		return
	}

	wsCtx, cancel := context.WithCancel(ctx)

	wsConn := core_ws_conn.NewConn(wsCtx, cancel, conn)

	wsConn.SetReadLimit(core_ws_client.MaxMessageSize)

	handler := core_ws_response.NewHandler(wsConn, logger, locale)

	limits, err := s.authLimiter.Allow(wsCtx, domain.RealtimeEventAuthRequest, ip)
	if err != nil {
		handler.AuthErrorResponse(err)

		_ = wsConn.Close()

		return
	}

	if !limits.Allowed {
		handler.AuthErrorResponse(core_errors.ErrTooManyRequests)

		_ = wsConn.Close()

		return
	}

	claims, err := s.AuthorizeClient(wsConn)
	if err != nil {
		handler.AuthErrorResponse(err)

		_ = wsConn.Close()

		return
	}

	client := core_ws_client.NewClient(
		wsConn,
		claims.UserID,
		claims.SessionID,
		s.presence,
	)

	if err := s.hub.Register(client); err != nil {
		handler.ErrorResponse(err, "failed to register websocket client")

		_ = wsConn.Close()

		return
	}

	client.Serve(s.hub.Unregister)
}
