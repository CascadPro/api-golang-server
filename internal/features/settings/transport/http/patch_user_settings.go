package settings_transport_http

import (
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	settings_http_dto "github.com/CascadePro/api-golang-server/internal/features/settings/transport/http/dto"
)

type PatchUserSettingsRequest settings_http_dto.PatchUserSettingsRequest

func (h *HttpHandler) PatchUserSettings(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID from context")
		return
	}

	var request PatchUserSettingsRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	settingsPatch := userSettingsPatchFromRequest(request)

	if _, err := h.settingsService.PatchUserSettings(ctx, userID, settingsPatch); err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user settings")
		return
	}

	responseHandler.NoContentResponse()
}

func userSettingsPatchFromRequest(request PatchUserSettingsRequest) domain.UserSettingsPatch {
	return domain.NewUserSettingsPatch(
		request.SessionExpireTerm.ToDomain(),
	)
}
