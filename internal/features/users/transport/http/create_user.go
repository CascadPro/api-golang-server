package users_transport_http

import (
	"net/http"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	users_http_dto "github.com/Svat-dev/golang-todo/internal/features/users/transport/http/dto"
)

type CreateUserRequest struct {
	FullName    string  `json:"full_name" validate:"required,min=2,max=100" example:"John Doe"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+" example:"+79999999999"`
}

type CreateUserResponse users_http_dto.UserDtoResponse

// CreateUser 	godoc
// @Summary 		Create a user
// @Description Create a new user in the system
// @Tags 				users
// @Accept 			json
// @Produce 		json
// @Param 			request body CreateUserRequest true "CreateUser body request"
// @Success 		200 {object} CreateUserResponse "Successfully created user"
// @Failure 		401 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/users [post]
func (h *UsersHttpHandler) CreateUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	var request CreateUserRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}

	userDomain := domainFromDto(request)

	userDomain, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")

		return
	}

	response := CreateUserResponse(users_http_dto.UserDtoFromDomain(userDomain))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func domainFromDto(dto CreateUserRequest) domain.User {
	return domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}
