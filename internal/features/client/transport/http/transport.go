package client_transport_http

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type HttpHandler struct {
	clientService ClientService
	tokenIssuer   core_jwt_security.IssuerMethods
	rateLimiter   core_http_middleware.HTTPRateLimitMiddleware
}

type ClientService interface {
	Get(context.Context, uuid.UUID) (domain.Client, error)

	Create(context.Context, domain.Client) (domain.Client, error)

	Delete(context.Context, uuid.UUID) error
}

func NewHttpHandler(
	clientService ClientService,
	tokenIssuer core_jwt_security.IssuerMethods,
	rateLimiter core_http_middleware.HTTPRateLimitMiddleware,
) *HttpHandler {
	return &HttpHandler{
		clientService: clientService,
		tokenIssuer:   tokenIssuer,
		rateLimiter:   rateLimiter,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	var (
		defaultMiddlewares = []core_http_middleware.Middleware{
			core_http_middleware.Authorization(h.tokenIssuer, domain.RoleAdmin, domain.RoleDirector, domain.RoleClerk),
		}
	)

	return []core_http_server.Route{
		{
			Method:     core_http.MethodPost,
			Path:       "/",
			Handler:    h.CreateClient,
			Middleware: append(defaultMiddlewares, h.rateLimiter.Middleware()),
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
