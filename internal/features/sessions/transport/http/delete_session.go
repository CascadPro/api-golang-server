package sessions_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

// DeleteSession godoc
// @Summary     Delete session
// @Description Delete a session by its ID
// @Tags        sessions
// @Param       id path string true "Session ID"
// @Success     204 "Session deleted"
// @Failure     400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure     401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure     404 {object} core_http_response.ErrorResponse "Session not found"
// @Failure     429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure     500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router      /sessions/{id} [delete]
func (h *HttpHandler) DeleteSession(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	sessionID, err := core_http_request.GetIDPathValue(r, "id", domain.SessionIDByteLength)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get session id path value")
		return
	}

	if err := h.sessionsService.DeleteSession(ctx, sessionID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete session")
		return
	}

	responseHandler.NoContentResponse()
}
