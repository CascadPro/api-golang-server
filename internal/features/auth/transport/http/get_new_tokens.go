package auth_transport_http

import (
	"net/http"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_http_utils "github.com/CascadePro/api-golang-server/internal/core/transport/http/utils"
)

const refreshTokenCookie = "refresh_token"

type GetNewTokensResponse struct {
	AccessToken string `json:"access_token" example:"JWT Token"`
}

// GetNewTokens godoc
// @Summary      Refresh access token
// @Description  Get new JWT‑access token using refresh token from cookie
// @Tags         auth
// @Produce      json
// @Param        refresh_token cookie string true "Refresh token"
// @Success      200 {object} GetNewTokensResponse "New access token"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /auth/login/refresh [get]
func (h *HttpHandler) GetNewTokens(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	rtCookie, err := core_http_utils.ParseCookie(r, refreshTokenCookie)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get refresh token from http request")

		return
	}

	accessToken, err := h.authService.GetNewTokens(ctx, rtCookie.Value)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get new tokens")

		return
	}

	if err := core_http_response.DeleteCookie(rw, "refresh_token"); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete cookie")

		return
	}

	response := GetNewTokensResponse{
		AccessToken: accessToken,
	}

	responseHandler.JsonResponse(response, http.StatusOK)
}
