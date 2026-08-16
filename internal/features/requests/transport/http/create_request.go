package requests_transport_http

import (
	"net/http"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	request_http_dto "github.com/CascadePro/api-golang-server/internal/features/requests/transport/http/dto"
)

type CreateRequestRequest struct {
	request_http_dto.CreateRequestRequest
}

type CreateRequestResponse request_http_dto.GetRequestsResponseRequest

// CreateRequest godoc
// @Summary      Create request
// @Description  Create request in the system
// @Tags         requests
// @Accept       json
// @Param 			 request body request_http_dto.CreateRequestRequest true "Create request body request"
// @Produce      json
// @Success      200 {object} CreateRequestResponse "Successfully created request"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /requests [post]
func (h *HttpHandler) CreateRequest(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	var request CreateRequestRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate http request")
		return
	}

	requestDomain := request_http_dto.CreateRequestDomainFromDTO(request.CreateRequestRequest)

	requestDomain, err := h.requestsService.CreateRequest(ctx, requestDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create request")
		return
	}

	response := CreateRequestResponse(request_http_dto.RequestsResponseFromDomain(requestDomain))

	responseHandler.JsonResponse(response, http.StatusOK)
}
