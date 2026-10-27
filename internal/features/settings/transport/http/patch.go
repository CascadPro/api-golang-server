package settings_transport_http

import (
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_http_types "github.com/CascadePro/api-golang-server/internal/core/transport/http/types"
	settings_http_dto "github.com/CascadePro/api-golang-server/internal/features/settings/transport/http/dto"
)

type PatchUserSettingsRequest struct {
	settings_http_dto.PatchUserSettingsRequest
}

// PatchUserSettings godoc
// @Summary      Patch user settings
// @Description  Atomically patch any user settings using optimistic locking
// @Tags         settings
// @Accept       json
// @Produce      json
// @Param        request body PatchUserSettingsRequest true "Partial user settings"
// @Success      200 {object} GetUserSettingsResponse "Updated user settings"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Settings not found"
// @Failure      409 {object} core_http_response.ErrorResponse "Settings version conflict"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /settings/user [patch]
func (h *HttpHandler) PatchUserSettings(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	locale := core_context.Locale(ctx)

	responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

	userID, err := core_context.UserID(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID from context")
		return
	}

	var request PatchUserSettingsRequest
	if err := core_http_request.DecodeAndValidate(rw, r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	patch := userSettingsPatchFromRequest(request.PatchUserSettingsRequest)
	settings, err := h.settingsService.PatchUserSettings(ctx, userID, patch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user settings")
		return
	}

	response := getUserSettingsDTOFromDomain(settings)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func nullablePointer[T any](value core_http_types.Nullable[T]) *T {
	if !value.Set || value.Value == nil {
		return nil
	}
	return value.Value
}

func userSettingsPatchFromRequest(request settings_http_dto.PatchUserSettingsRequest) domain.UserSettingsPatch {
	return domain.NewUserSettingsPatch(
		request.Version,
		request.SessionExpireTerm.ToDomain(),
		request.Session2FAEnabled.ToDomain(),
		domain.UserPreferencesPatch{
			Theme:                   nullablePointer(request.Preferences.Theme),
			Accent:                  nullablePointer(request.Preferences.Accent),
			Animations:              nullablePointer(request.Preferences.Animations),
			Density:                 nullablePointer(request.Preferences.Density),
			TextSize:                nullablePointer(request.Preferences.TextSize),
			ConfirmImportantActions: nullablePointer(request.Preferences.ConfirmImportantActions),
		},
		domain.UserNotificationPreferencesPatch{
			Requests: nullablePointer(request.Notifications.Requests),
			Sessions: nullablePointer(request.Notifications.Sessions),
			System:   nullablePointer(request.Notifications.System),
			Push:     nullablePointer(request.Notifications.Push),
		},
	)
}
