package root_transport_http

import (
	"net/http"

	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
)

type HttpHandler struct{}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodGet,
			Path:    "/",
			Handler: h.Root,
		},
		{
			Method:  http.MethodGet,
			Path:    "/ping",
			Handler: h.Ping,
		},
	}
}
