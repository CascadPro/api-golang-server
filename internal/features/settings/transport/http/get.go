package settings_transport_http

import (
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	settings_http_dto "github.com/CascadePro/api-golang-server/internal/features/settings/transport/http/dto"
)

type GetUserSettingsResponse settings_http_dto.UserSettingsResponse

// GetUserSettings godoc
// @Summary      Get user settings
// @Description  Get all settings for the authenticated user
// @Tags         settings
// @Produce      json
// @Success      200 {object} GetUserSettingsResponse "User settings"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Not found"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /settings/user [get]
func (h *HttpHandler) GetUserSettings(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	locale := core_context.Locale(ctx)

	responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

	userID, err := core_context.UserID(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID from context")
		return
	}

	settings, err := h.settingsService.GetUserSettings(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user settings")
		return
	}

	response := getUserSettingsDTOFromDomain(settings)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func getUserSettingsDTOFromDomain(settings domain.UserSettings) GetUserSettingsResponse {
	return GetUserSettingsResponse{
		ID: settings.ID, Version: settings.Version,
		SessionExpireTerm: settings.SessionExpireTerm,
		Session2FAEnabled: settings.Session2FAEnabled,
		Preferences: settings_http_dto.UserSettingsResponsePreferences{
			Theme:                   settings.Preferences.Theme,
			Accent:                  settings.Preferences.Accent,
			Animations:              settings.Preferences.Animations,
			Density:                 settings.Preferences.Density,
			TextSize:                settings.Preferences.TextSize,
			ConfirmImportantActions: settings.Preferences.ConfirmImportantActions,
		},
		Notifications: settings_http_dto.UserSettingsResponseNotifications{
			Requests: settings.Notifications.Requests,
			Sessions: settings.Notifications.Sessions,
			System:   settings.Notifications.System,
			Push:     settings.Notifications.Push,
		},
	}
}
