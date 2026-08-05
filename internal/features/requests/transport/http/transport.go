package requests_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
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

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
	}
}
