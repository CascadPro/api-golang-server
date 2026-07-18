package users_transport_http

import (
	"net/http"

	core_logger "github.com/Svat-dev/golang-todo/internal/core/logger"
	core_http_request "github.com/Svat-dev/golang-todo/internal/core/transport/http/request"
	core_http_response "github.com/Svat-dev/golang-todo/internal/core/transport/http/response"
)

// DeleteUser		godoc
// @Summary 		Delete user
// @Description Delete a user in the system
// @Tags 				users
// @Param       id path int true "User ID"
// @Success 		204 "No content response"
// @Failure 		401 {object} core_http_response.ErrorResponse "Bad request"
// @Failure 		409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure		 	500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router		 	/users/{id} [delete]
func (h *UsersHttpHandler) DeleteUser(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, rw)

	userID, err := core_http_request.GetIntPathValue(r, "id")
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
