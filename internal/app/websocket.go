package app

import (
	"context"
	"fmt"

	core_ws_dispatcher "github.com/CascadePro/api-golang-server/internal/core/transport/ws/dispatcher"
	core_ws_hub "github.com/CascadePro/api-golang-server/internal/core/transport/ws/hub"
	core_ws_server "github.com/CascadePro/api-golang-server/internal/core/transport/ws/server"
)

func (a *App) initWebsocket(ctx context.Context) error {
	hub := core_ws_hub.NewHub(ctx, a.infrastructure.Publisher, a.features.Presence, a.logger)
	if err := a.infrastructure.Publisher.Subscribe(ctx, hub.HandleEvent); err != nil {
		return fmt.Errorf("initialize websocket subscriber: %w", err)
	}

	authRateLimiter, err := a.newWSRateLimiter(a.rateLimits.WebSocketAuth)
	if err != nil {
		return fmt.Errorf("initialize websocket auth rate limiter: %w", err)
	}

	dispatcher := core_ws_dispatcher.NewDispatcher(a.features.Presence, a.infrastructure.TokenIssuer)

	wsServer := core_ws_server.NewServer(hub, dispatcher, a.infrastructure.WsAuthenticator, authRateLimiter, a.cfg.AllowedOrigins)

	a.wsServer = wsServer

	return nil
}
