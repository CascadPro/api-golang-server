package users_transport_http

import (
	"net/http"

	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	core_http_utils "github.com/Svat-dev/golang-todo/internal/core/transport/http/utils"
	users_http_dto "github.com/Svat-dev/golang-todo/internal/features/users/transport/http/dto"
)

type GetUserResponse users_http_dto.UserDtoResponse

func (h *UsersHttpHandler) GetUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")

		return
	}

	user, err := h.usersService.GetUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user")

		return
	}

	response := users_http_dto.UserDtoResponse(users_http_dto.UserDtoFromDomain(user))

	responseHandler.JsonResponse(response, http.StatusOK)
}
