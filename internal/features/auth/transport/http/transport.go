package auth_transport_http

import (
	"net/http"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	auth_service "github.com/CascadePro/api-golang-server/internal/features/auth/service"
)

type HttpHandler struct {
	authService auth_service.ServiceMethods
}

func NewHttpHandler(authService auth_service.ServiceMethods) *HttpHandler {
	return &HttpHandler{
		authService: authService,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	rateLimit := core_http_middleware.NewRateLimitConfig(5, 5*time.Minute)
	roles := []domain.UserRole{domain.RoleAdmin, domain.RoleDirector}

	return []core_http_server.Route{
		{
			Method:     http.MethodPost,
			Path:       "/register/token",
			Handler:    h.CreateRegisterToken,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Authorization(roles...)},
		},
		{
			Method:     http.MethodPost,
			Path:       "/register",
			Handler:    h.Register,
			Middleware: []core_http_middleware.Middleware{rateLimit.Middleware()},
		},
		{
			Method:     http.MethodPost,
			Path:       "/login",
			Handler:    h.Login,
			Middleware: []core_http_middleware.Middleware{rateLimit.Middleware()},
		},
		{
			Method:     http.MethodPost,
			Path:       "/logout",
			Handler:    h.Logout,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Authorization()},
		},
		{
			Method:  http.MethodGet,
			Path:    "/login/refresh",
			Handler: h.GetNewTokens,
		},
	}
}
