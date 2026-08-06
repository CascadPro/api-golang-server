package requests_transport_http

import (
	"net/http"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	request_http_dto "github.com/CascadePro/api-golang-server/internal/features/requests/transport/http/dto"
)

type GetRequestResponse request_http_dto.GetRequestResponse

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
