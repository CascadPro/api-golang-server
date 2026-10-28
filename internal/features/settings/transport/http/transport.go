package settings_transport_http

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type HttpHandler struct {
	settingsService SettingsService
}

type SettingsService interface {
	GetUserSettings(ctx context.Context, userID uuid.UUID) (domain.UserSettings, error)

	PatchUserSettings(ctx context.Context, userID uuid.UUID, patch domain.UserSettingsPatch) (domain.UserSettings, error)
}

func NewHttpHandler(settingsService SettingsService) *HttpHandler {
	return &HttpHandler{settingsService: settingsService}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  core_http.MethodGet,
			Path:    "/user",
			Handler: h.GetUserSettings,
		},
		{
			Method:  core_http.MethodPatch,
			Path:    "/user",
			Handler: h.PatchUserSettings,
		},
	}
}
