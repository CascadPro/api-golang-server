package sessions_transport_http

import (
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	session_service "github.com/CascadePro/api-golang-server/internal/features/sessions/service"
)

type HttpHandler struct {
	sessionsService session_service.ServiceMethods
}

func NewHttpHandler(sessionsService session_service.ServiceMethods) *HttpHandler {
	return &HttpHandler{
		sessionsService: sessionsService,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  core_http.MethodGet,
			Path:    "/",
			Handler: h.GetUserSessions,
		},
		{
			Method:  core_http.MethodDelete,
			Path:    "/{id}",
			Handler: h.DeleteSession,
		},
		{
			Method:  core_http.MethodDelete,
			Path:    "/delete",
			Handler: h.DeleteUserSessions,
		},
	}
}
