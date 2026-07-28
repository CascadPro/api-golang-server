package auth_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type RegisterRequest struct {
	Email    string `json:"email"    validate:"required,email" example:"test@example.com"`
	Password string `json:"password" validate:"required,pwd"   example:"Strong_Pwd123"`
	Token    string `json:"token"    validate:"required,uuid4" example:"00000000-000000-000000-000000000000"`
}

// Register godoc
// @Summary 		Register
// @Description User registration by "register token" from CreateRegisterToken
// @Tags 				auth
// @Accept 			json
// @Param 			request body RegisterRequest true "Register body request"
// @Success 		204 "Successfully created register token"
// @Failure 		400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure 		429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/auth/register [post]
func (h *HttpHandler) Register(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	var request RegisterRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}

	patch := patchFromRegisterRequest(request)

	if err := h.authService.Register(ctx, patch, request.Token); err != nil {
		responseHandler.ErrorResponse(err, "failed to register")

		return
	}

	responseHandler.NoContentResponse()
}

func patchFromRegisterRequest(request RegisterRequest) domain.UserPatch {
	return domain.NewUserRegisterPatch(
		domain.NewNullable(request.Email),
		domain.NewNullable(request.Password),
	)
}
