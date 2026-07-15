package users_transport_http

import (
	"net/http"

	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
	core_http_utils "github.com/Svat-dev/golang-todo/internal/core/transport/http/utils"
)

func (h *UsersHttpHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	userID, err := core_http_utils.GetIntPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID path value")

		return
	}

	err = h.usersService.DeleteUser(ctx, userID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user")

		return
	}

	responseHandler.NoContentResponse()
}
