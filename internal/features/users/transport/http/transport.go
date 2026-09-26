package users_transport_http

import (
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	users_service "github.com/CascadePro/api-golang-server/internal/features/users/service"
)

type HttpHandler struct {
	usersService users_service.ServiceMethods
	rateLimiter  core_http_middleware.HTTPRateLimitMiddleware
}

func NewHttpHandler(usersService users_service.ServiceMethods, rateLimiter core_http_middleware.HTTPRateLimitMiddleware) *HttpHandler {
	return &HttpHandler{
		usersService: usersService,
		rateLimiter:  rateLimiter,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  core_http.MethodGet,
			Path:    "/my",
			Handler: h.GetCurrentUser,
		},
		{
			Method:     core_http.MethodPatch,
			Path:       "/avatar",
			Handler:    h.UpdateAvatar,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Media(), h.rateLimiter.Middleware()},
		},
		{
			Method:     core_http.MethodDelete,
			Path:       "/avatar",
			Handler:    h.DeleteAvatar,
			Middleware: []core_http_middleware.Middleware{h.rateLimiter.Middleware()},
		},
	}
}
