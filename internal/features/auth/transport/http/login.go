package auth_transport_http

import (
	"fmt"
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type LoginRequest struct {
	Email    string `json:"email,omitempty"    validate:"omitempty,email" example:"test@example.com" swaggertype:"string"`
	Phone    string `json:"phone,omitempty"    validate:"omitempty"       example:"+79001234567" swaggertype:"string"`
	Password string `json:"password"           validate:"required,pwd"    example:"Strong_Pwd123"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"  example:"JWT Token"`
	RefreshToken string `json:"refresh_token" example:"JWT Token"`
}

// Login godoc
// @Summary 		Login
// @Description User login in the system using email or phone
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
	locale := core_context.Locale(ctx)

	responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

	var request LoginRequest
	if err := core_http_request.DecodeAndValidate(rw, r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}

	// Validate that either email or phone is provided
	if request.Email == "" && request.Phone == "" {
		responseHandler.ErrorResponse(fmt.Errorf("either email or phone must be provided"), "validation failed")
		return
	}
	if request.Email != "" && request.Phone != "" {
		responseHandler.ErrorResponse(fmt.Errorf("cannot provide both email and phone"), "validation failed")
		return
	}

	// Validate phone if provided
	if request.Phone != "" {
		if _, err := domain.ParsePhoneNumber(request.Phone); err != nil {
			responseHandler.ErrorResponse(fmt.Errorf("invalid phone number: %w", err), "validation failed")
			return
		}
	}

	ip, err := core_context.ClientIP(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get client ip")

		return
	}

	userAgent, err := core_http_request.ParseUserAgent(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user agent")

		return
	}

	var user domain.User
	if request.Email != "" {
		user = domain.NewUserLogin(request.Email, request.Password)
	} else {
		// Login with phone - we need to look up the user by phone first
		// This will be handled by the auth service
		p, _ := domain.ParsePhoneNumber(request.Phone)
		user = domain.User{
			ID:           domain.UninitializedUUID,
			Version:      domain.UninitializedVersion,
			Activated:    true,
			Role:         domain.UninitializedRole,
			Phone:        &p,
			PasswordHash: &request.Password,
		}
	}

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