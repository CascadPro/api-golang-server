package sessions_transport_http

import (
	"net/http"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

// DeleteUserSessions godoc
// @Summary     Delete user sessions
// @Description Delete all sessions connected with authorized user
// @Tags        sessions
// @Success     204 "Sessions deleted"
// @Failure     401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure     404 {object} core_http_response.ErrorResponse "Session not found"
// @Failure     429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /sessions/delete [delete]
func (h *HttpHandler) DeleteUserSessions(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	if err := h.sessionsService.DeleteUserSessions(ctx); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete user session")
		return
	}

	responseHandler.NoContentResponse()
}
