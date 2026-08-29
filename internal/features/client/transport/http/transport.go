package client_transport_http

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	client_service "github.com/CascadePro/api-golang-server/internal/features/client/service"
)

type HttpHandler struct {
	clientService client_service.ServiceMethods
	tokenIssuer   core_jwt_security.IssuerMethods
}

func NewHttpHandler(clientService client_service.ServiceMethods, tokenIssuer core_jwt_security.IssuerMethods) *HttpHandler {
	return &HttpHandler{
		clientService: clientService,
		tokenIssuer:   tokenIssuer,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	var (
		defaultMiddlewares = []core_http_middleware.Middleware{
			core_http_middleware.Authorization(h.tokenIssuer, domain.RoleAdmin, domain.RoleDirector, domain.RoleClerk),
		}
		rateLimitCfg = core_http_middleware.NewRateLimitConfig(15, 5*time.Minute)
	)

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
