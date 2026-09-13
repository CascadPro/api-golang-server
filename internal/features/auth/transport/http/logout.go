package auth_transport_http

import (
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_http_utils "github.com/CascadePro/api-golang-server/internal/core/transport/http/utils"
)

// Logout godoc
// @Summary 		Logout
// @Description User logout
// @Tags 				auth
// @Success 		204 "Successfully logged out"
// @Failure 		400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure     429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure 		500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/auth/logout [post]
func (h *HttpHandler) Logout(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	locale := core_context.Locale(ctx)

	responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

	cookie, err := core_http_utils.ParseCookie(r, refreshTokenCookie)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get refresh token from http request")
		return
	}

	if err := h.authService.Logout(ctx, cookie.Value); err != nil {
		responseHandler.ErrorResponse(err, "failed to logout")
		return
	}

	_ = core_http_response.DeleteCookie(rw, refreshTokenCookie)

	responseHandler.NoContentResponse()
}
