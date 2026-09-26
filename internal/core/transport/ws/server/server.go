package core_ws_server

import (
	"context"
	"net/http"
	"strings"

	core_ws_hub "github.com/CascadePro/api-golang-server/internal/core/transport/ws/hub"
	core_ws_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/ws/middleware"
	presence_service "github.com/CascadePro/api-golang-server/internal/features/presence/service"
	"github.com/gorilla/websocket"
)

type Server struct {
	hub      *core_ws_hub.Hub
	presence presence_service.ServiceMethods

	authenticator core_ws_middleware.AuthenticatorMethods
	authLimiter   *core_ws_middleware.WSRateLimiter

	upgrader websocket.Upgrader
}

func NewServer(
	hub *core_ws_hub.Hub,
	presence presence_service.ServiceMethods,
	authenticator core_ws_middleware.AuthenticatorMethods,
	authLimiter *core_ws_middleware.WSRateLimiter,
	allowedOrigins string,
) *Server {
	origins := make(map[string]struct{})

	for _, origin := range strings.Split(allowedOrigins, ",") {
		origin = strings.TrimSpace(origin)

		if origin != "" {
			origins[origin] = struct{}{}
		}
	}

	return &Server{
		hub:      hub,
		presence: presence,

		authenticator: authenticator,
		authLimiter:   authLimiter,

		upgrader: websocket.Upgrader{
			ReadBufferSize:  4096,
			WriteBufferSize: 4096,

			CheckOrigin: func(r *http.Request) bool {
				origin := r.Header.Get("Origin")

				if origin == "" {
					return true
				}

				_, ok := origins[origin]

				return ok
			},
		},
	}
}

func (s *Server) Shutdown(ctx context.Context) error {
	return s.hub.Shutdown(ctx)
}

var _ http.Handler = (*Server)(nil)
