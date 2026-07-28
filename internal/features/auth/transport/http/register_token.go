package auth_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type CreateRegisterTokenRequest struct {
	Name     string          `json:"name"      validate:"required,name"      example:"John"`
	Surname  string          `json:"surname"   validate:"required,name"      example:"Doe"`
	LastName *string         `json:"last_name" validate:"omitempty,name"     example:""`
	Role     domain.UserRole `json:"role"      validate:"required,user_role" example:"regular"`
}

type CreateRegisterTokenResponse struct {
	Token string `json:"token" example:"00000000-000000-000000-000000000000"`
}

// CreateRegisterToken godoc
// @Summary 		Create register token
// @Description Create new user and `token` to continue registration
// @Tags 				auth
// @Accept 			json
// @Produce 		json
// @Param 			request body CreateRegisterTokenRequest true "CreateRegisterToken body request"
// @Success 		200 {object} CreateRegisterTokenResponse "Successfully created register token"
// @Failure 		400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure 		403 {object} core_http_response.ErrorResponse "Forbidden"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure 		429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/auth/register/token [post]
func (h *HttpHandler) CreateRegisterToken(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	var request CreateRegisterTokenRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}

	user := domain.NewRegisterUser(request.Name, request.Surname, request.LastName)

	token, err := h.authService.CreateRegisterToken(ctx, user)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create register token")

		return
	}

	response := CreateRegisterTokenResponse{
		Token: token.Token,
	}

	responseHandler.JsonResponse(response, http.StatusOK)
}
