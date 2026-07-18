package users_transport_http

import (
	"net/http"

	"github.com/Svat-dev/golang-todo/internal/core/domain"
	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	users_http_dto "github.com/Svat-dev/golang-todo/internal/features/users/transport/http/dto"
)

type PatchUserRequest users_http_dto.PatchUserRequest

type PatchUserResponse users_http_dto.UserDtoResponse

// PatchUser 		godoc
// @Summary 		Patch a user
// @Description Patch a user in the system
// @Tags 				users
// @Accept 			json
// @Produce 		json
// @Param       id path int true "User ID"
// @Param 			request body users_http_dto.PatchUserRequestSwagger true "PatchUser body request"
// @Success 		200 {object} PatchUserResponse "Successfully patched user"
// @Failure 		401 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/users/{id} [patch]
func (h *UsersHttpHandler) PatchUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")

		return
	}

	var request PatchUserRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")

		return
	}

	userPatch := userPatchFromRequest(request)

	userDomain, err := h.usersService.PatchUser(ctx, userID, userPatch)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to patch user")

		return
	}

	response := PatchUserResponse(users_http_dto.UserDtoFromDomain(userDomain))

	responseHandler.JsonResponse(response, http.StatusOK)
}

func userPatchFromRequest(request PatchUserRequest) domain.UserPatch {
	return domain.NewUserPatch(
		request.FullName.ToDomain(),
		request.PhoneNumber.ToDomain(),
	)
}
