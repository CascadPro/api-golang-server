package client_transport_http

import (
	"net/http"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

// DeleteClient godoc
// @Summary      Delete client
// @Description  Delete client from the system
// @Tags         clients
// @Param        id path string true "Client ID"
// @Success      204 "Successfully deleted client"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /clients/{id} [delete]
func (h *HttpHandler) DeleteClient(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	requestID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request id from path value")
		return
	}

	if err := h.clientService.DeleteClient(ctx, requestID); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete client")
		return
	}

	responseHandler.NoContentResponse()
}
