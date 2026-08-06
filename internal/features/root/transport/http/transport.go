package root_transport_http

import (
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
)

type HttpHandler struct{}

func NewHttpHandler() *HttpHandler {
	return &HttpHandler{}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  core_http.MethodGet,
			Path:    "/",
			Handler: h.Root,
		},
		{
			Method:  core_http.MethodGet,
			Path:    "/ping",
			Handler: h.Ping,
		},
	}
}
