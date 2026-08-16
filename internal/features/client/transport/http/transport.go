package client_transport_http

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	client_service "github.com/CascadePro/api-golang-server/internal/features/client/service"
)

type HttpHandler struct {
	clientService client_service.ServiceMethods
}

func NewHttpHandler(clientService client_service.ServiceMethods) *HttpHandler {
	return &HttpHandler{
		clientService: clientService,
	}
}

var (
	defaultMiddlewares = []core_http_middleware.Middleware{
		core_http_middleware.Authorization(domain.RoleAdmin, domain.RoleDirector, domain.RoleClerk),
	}
)

func (h *HttpHandler) Routes() []core_http_server.Route {
	rateLimitCfg := core_http_middleware.NewRateLimitConfig(15, 5*time.Minute)

	return []core_http_server.Route{
		{
			Method:     core_http.MethodPost,
			Path:       "/",
			Handler:    h.CreateClient,
			Middleware: append(defaultMiddlewares, rateLimitCfg.Middleware()),
		},
		{
			Method:     core_http.MethodGet,
			Path:       "/{id}",
			Handler:    h.GetClient,
			Middleware: defaultMiddlewares,
		},
		{
			Method:     core_http.MethodDelete,
			Path:       "/{id}",
			Handler:    h.DeleteClient,
			Middleware: defaultMiddlewares,
		},
	}
}
