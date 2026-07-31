package users_transport_http

import (
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	users_service "github.com/CascadePro/api-golang-server/internal/features/users/service"
)

type HttpHandler struct {
	usersService users_service.ServiceMethods
}

func NewHttpHandler(usersService users_service.ServiceMethods) *HttpHandler {
	return &HttpHandler{
		usersService: usersService,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{}
}
