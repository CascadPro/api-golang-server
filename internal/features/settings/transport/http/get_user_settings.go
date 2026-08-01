package settings_transport_http

import (
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	"github.com/google/uuid"
)

type GetUserSettingsResponse struct {
	ID                uuid.UUID                `json:"id"                  example:"00000000-000000-000000-000000000000"`
	Version           int                      `json:"version"             example:"1"`
	SessionExpireTerm domain.SessionExpireTime `json:"session_expire_term" example:"30d"`
}

func (h *HttpHandler) GetUserSettings(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID from context")
		return
	}

	settings, err := h.settingsService.GetUserSettings(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user settings")
		return
	}

	response := getUserSettingsDtoFromDomain(settings)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func getUserSettingsDtoFromDomain(settings domain.UserSettings) GetUserSettingsResponse {
	return GetUserSettingsResponse{
		ID:                settings.ID,
		Version:           settings.Version,
		SessionExpireTerm: settings.SessionExpireTerm,
	}
}
