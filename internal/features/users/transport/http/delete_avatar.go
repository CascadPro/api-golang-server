package users_transport_http

import (
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

// DeleteAvatar godoc
// @Summary      Delete avatar
// @Description  Delete user's avatar for the authenticated user
// @Tags         users
// @Success      204 "Successfully deleted user avatar"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Not found"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /users/avatar [delete]
func (h *HttpHandler) DeleteAvatar(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get userID from context")
		return
	}

	if err := h.usersService.DeleteAvatar(ctx, userID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete avatar")
		return
	}

	responseHandler.NoContentResponse()
}
