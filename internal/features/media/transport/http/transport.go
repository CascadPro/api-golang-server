package media_transport_http

import (
	"net/http"

	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
)

type HttpHandler struct {
	mediaService media_service.ServiceMethods
}

func NewHttpHandler(mediaService media_service.ServiceMethods) *HttpHandler {
	return &HttpHandler{
		mediaService: mediaService,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/{tag}/{filename}",
			Handler: h.GetFile,
		},
	}
}
