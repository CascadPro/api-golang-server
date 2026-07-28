package sessions_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	sessions_http_dto "github.com/CascadePro/api-golang-server/internal/features/sessions/transport/http/dto"
)

type GetUserSessionsResponse struct {
	Current  sessions_http_dto.SessionDto   `json:"current_session"`
	Sessions []sessions_http_dto.SessionDto `json:"sessions"`
}

// GetUserSessions godoc
// @Summary      Get user sessions
// @Description  Get all sessions for the authenticated user
// @Tags         sessions
// @Produce      json
// @Success      200 {object} GetUserSessionsResponse "User sessions"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "User sessions not found"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /sessions [get]
func (h *HttpHandler) GetUserSessions(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	sessions, err := h.sessionsService.GetUserSessions(ctx)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get user sessions")

		return
	}

	response := sessionsResponseFromDomain(sessions)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func sessionsResponseFromDomain(sessions []domain.Session) GetUserSessionsResponse {
	return GetUserSessionsResponse{
		Current:  sessions_http_dto.SessionDomainToDTO(sessions[0]),
		Sessions: sessions_http_dto.SessionDomainsToDTOs(sessions[1:]),
	}
}
