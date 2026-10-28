package requests_transport_http

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
	requestsService RequestsService
	tokenIssuer     core_jwt_security.IssuerMethods
	rateLimiter     core_http_middleware.HTTPRateLimitMiddleware
}

type RequestsService interface {
	CreateRequest(context.Context, domain.Request) (domain.Request, error)

	GetRequest(context.Context, uuid.UUID) (domain.Request, domain.User, domain.Client, map[int]domain.File, error)
	GetRequests(context.Context, int, int, domain.RequestQueryFilter) ([]domain.Request, int64, error)
	GetRequestsStats(context.Context) (domain.RequestStats, error)

	PatchRequest(context.Context, uuid.UUID, domain.RequestPatch) (domain.Request, error)
	PatchRequestStatus(context.Context, uuid.UUID, domain.RequestStatus) error

	UploadDoc(context.Context, uuid.UUID, domain.File, []byte, int) error
	DeleteDoc(context.Context, uuid.UUID, int) error
}

func NewHttpHandler(
	requestsService RequestsService,
	tokenIssuer core_jwt_security.IssuerMethods,
	rateLimiter core_http_middleware.HTTPRateLimitMiddleware,
) *HttpHandler {
	return &HttpHandler{
		requestsService: requestsService,
		tokenIssuer:     tokenIssuer,
		rateLimiter:     rateLimiter,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	var (
		adminMiddlewares = []core_http_middleware.Middleware{
			core_http_middleware.Authorization(h.tokenIssuer, domain.RoleAdmin, domain.RoleDirector),
		}
		defaultMiddlewares = []core_http_middleware.Middleware{
			core_http_middleware.Authorization(h.tokenIssuer, domain.RoleAdmin, domain.RoleDirector, domain.RoleClerk),
		}
	)

	return []core_http_server.Route{
		{
			Method:     core_http.MethodQuery,
			Path:       "/",
			Handler:    h.GetRequests,
			Middleware: defaultMiddlewares,
		},
		{
			Method:     core_http.MethodPost,
			Path:       "/",
			Handler:    h.CreateRequest,
			Middleware: defaultMiddlewares,
		},
		{
			Method:     core_http.MethodGet,
			Path:       "/stats",
			Handler:    h.GetRequestStats,
			Middleware: defaultMiddlewares,
		},
		{
			Method:     core_http.MethodGet,
			Path:       "/{id}",
			Handler:    h.GetRequest,
			Middleware: defaultMiddlewares,
		},
		{
			Method:     core_http.MethodPatch,
			Path:       "/{id}/reject",
			Handler:    h.RejectRequest,
			Middleware: adminMiddlewares,
		},
		{
			Method:     core_http.MethodPatch,
			Path:       "/{id}/approve",
			Handler:    h.ApproveRequest,
			Middleware: adminMiddlewares,
		},
		{
			Method:     core_http.MethodPatch,
			Path:       "/{id}",
			Handler:    h.PatchRequest,
			Middleware: defaultMiddlewares,
		},
		{
			Method:     core_http.MethodPost,
			Path:       "/{id}/file/{index}",
			Handler:    h.UploadDoc,
			Middleware: append(defaultMiddlewares, h.rateLimiter.Middleware(), core_http_middleware.Media()),
		},
		{
			Method:     core_http.MethodDelete,
			Path:       "/{id}/file/{index}",
			Handler:    h.DeleteDoc,
			Middleware: append(defaultMiddlewares, h.rateLimiter.Middleware()),
		},
	}
}
