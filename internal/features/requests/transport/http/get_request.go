package requests_transport_http

import (
	"net/http"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	request_http_dto "github.com/CascadePro/api-golang-server/internal/features/requests/transport/http/dto"
)

type GetRequestResponse request_http_dto.GetRequestResponse

// GetRequest godoc
// @Summary      Get request
// @Description  Get request by ID
// @Tags         requests
// @Param        id path string true "Request ID"
// @Produce      json
// @Success      200 {object} GetRequestResponse "Request"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Request not found"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /requests/{id} [get]
func (h *HttpHandler) GetRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	requestID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request id path value")
		return
	}

	request, user, files, err := h.requestsService.GetRequest(ctx, requestID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request")
		return
	}

	response := request_http_dto.RequestResponseFromDomain(request, user, &files)

	responseHandler.JsonResponse(response, http.StatusOK)
}
