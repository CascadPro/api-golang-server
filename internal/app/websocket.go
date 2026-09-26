package app

import (
	"context"
	"fmt"
	"time"

	core_ws_dispatcher "github.com/CascadePro/api-golang-server/internal/core/transport/ws/dispatcher"
	core_ws_hub "github.com/CascadePro/api-golang-server/internal/core/transport/ws/hub"
	core_ws_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/ws/middleware"
	core_ws_server "github.com/CascadePro/api-golang-server/internal/core/transport/ws/server"
)

const (
	authRateLimit = 10
)

func (a *App) initWebsocket(ctx context.Context) error {
	hub := core_ws_hub.NewHub(
		ctx,
		a.infrastructure.Publisher,
		a.features.Presence,
		a.logger,
	)

	if err := a.infrastructure.Publisher.Subscribe(ctx, hub.HandleEvent); err != nil {
		return fmt.Errorf("initialize websocket subscriber: %w", err)
	}

	authRateLimiter := core_ws_middleware.NewWSRateLimiter(
		a.infrastructure.Redis,
		authRateLimit,
		time.Minute,
	)

	dispatcher := core_ws_dispatcher.NewDispatcher(
		a.features.Presence,
		a.infrastructure.TokenIssuer,
	)

	wsServer := core_ws_server.NewServer(
		hub,
		dispatcher,
		a.infrastructure.WsAuthenticator,
		authRateLimiter,
		a.cfg.AllowedOrigins,
	)

	a.wsServer = wsServer

	return nil
}
