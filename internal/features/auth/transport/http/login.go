package auth_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_http_utils "github.com/CascadePro/api-golang-server/internal/core/transport/http/utils"
)

type LoginRequest struct {
	Email    string `json:"email"    validate:"required,email" example:"test@example.com"`
	Password string `json:"password" validate:"required,pwd"   example:"Strong_Pwd123"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"  example:"JWT Token"`
	RefreshToken string `json:"refresh_token" example:"JWT Token"`
}

// Login godoc
// @Summary 		Login
// @Description User login in the system
// @Tags 				auth
// @Accept 			json
// @Produce 		json
// @Param 			request body LoginRequest true "Login body request"
// @Success 		200 {object} LoginResponse "Two JWT tokens: access and refresh"
// @Failure 		400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 		404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure 		429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/auth/login [post]
func (h *HttpHandler) Login(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}

	ip, err := core_http_utils.ClientIP(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get client ip")

		return
	}

	userAgent, err := core_http_request.ParseUserAgent(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user agent")

		return
	}

	user := domain.NewUserLogin(request.Email, request.Password)

	accessToken, refreshToken, err := h.authService.Login(ctx, user, ip, userAgent)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to login")

		return
	}

	response := LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	responseHandler.JsonResponse(response, http.StatusOK)
}
