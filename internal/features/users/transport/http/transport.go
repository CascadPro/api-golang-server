package users_transport_http

import (
	"net/http"

	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
	users_service "github.com/CascadePro/api-golang-server/internal/features/users/service"
)

type HttpHandler struct {
	usersService users_service.ServiceMethods
	mediaService media_service.ServiceMethods
}

func NewHttpHandler(usersService users_service.ServiceMethods, mediaService media_service.ServiceMethods) *HttpHandler {
	return &HttpHandler{
		usersService: usersService,
		mediaService: mediaService,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/my",
			Handler: h.GetCurrentUser,
		},
	}
}
