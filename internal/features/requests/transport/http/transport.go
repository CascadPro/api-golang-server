package requests_transport_http

import (
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	requests_service "github.com/CascadePro/api-golang-server/internal/features/requests/service"
)

type HttpHandler struct {
	requestsService requests_service.ServiceMethods
}

func NewHttpHandler(requestsService requests_service.ServiceMethods) *HttpHandler {
	return &HttpHandler{
		requestsService: requestsService,
	}
}

var (
	adminMiddlewares = []core_http_middleware.Middleware{
		core_http_middleware.Authorization(domain.RoleAdmin, domain.RoleDirector),
	}
	defaultMiddlewares = []core_http_middleware.Middleware{
		core_http_middleware.Authorization(domain.RoleAdmin, domain.RoleDirector, domain.RoleClerk),
	}
)

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:     core_http.MethodQuery,
			Path:       "/",
			Handler:    h.GetRequests,
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
			Path:       "/{id}/update",
			Handler:    h.PatchRequest,
			Middleware: defaultMiddlewares,
		},
		{
			Method:     core_http.MethodPost,
			Path:       "/{id}/upload",
			Handler:    h.UploadDoc,
			Middleware: append(adminMiddlewares, core_http_middleware.Media()),
		},
	}
}
