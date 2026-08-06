package requests_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

func (h *HttpHandler) ApproveRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	requestID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request id path value")
		return
	}

	if err := h.requestsService.PatchRequestStatus(ctx, requestID, domain.RequestStatusApproved); err != nil {
		responseHandler.ErrorResponse(err, "failed to approve request")
		return
	}

	responseHandler.NoContentResponse()
}
