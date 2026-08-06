package requests_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

// RejectRequest godoc
// @Summary      Reject a request
// @Description  Reject a request by ID
// @Tags         requests
// @Param        id path string true "Request ID"
// @Success      204 "Rejected request"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Request not found"
// @Failure 		 409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /requests/{id}/reject [patch]
func (h *HttpHandler) RejectRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	requestID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request id path value")
		return
	}

	if err := h.requestsService.PatchRequestStatus(ctx, requestID, domain.RequestStatusRejected); err != nil {
		responseHandler.ErrorResponse(err, "failed to reject request")
		return
	}

	responseHandler.NoContentResponse()
}
