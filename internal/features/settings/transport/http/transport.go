package settings_transport_http

import (
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	settings_service "github.com/CascadePro/api-golang-server/internal/features/settings/service"
)

type HttpHandler struct {
	settingsService settings_service.ServiceMethods
}

func NewHttpHandler(settingsService settings_service.ServiceMethods) *HttpHandler {
	return &HttpHandler{
		settingsService: settingsService,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  core_http.MethodGet,
			Path:    "/user/my",
			Handler: h.GetUserSettings,
		},
		{
			Method:  core_http.MethodPatch,
			Path:    "/user/update",
			Handler: h.PatchUserSettings,
		},
	}
}
