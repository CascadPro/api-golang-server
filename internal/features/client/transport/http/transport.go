package client_transport_http

import (
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

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{}
}
