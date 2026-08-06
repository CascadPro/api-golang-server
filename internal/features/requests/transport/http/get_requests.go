package requests_transport_http

import (
	"errors"
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	request_http_dto "github.com/CascadePro/api-golang-server/internal/features/requests/transport/http/dto"
)

type GetRequestsRequest request_http_dto.GetRequestsRequest

type GetRequestsResponse request_http_dto.GetRequestsResponse

// GetRequests godoc
// @Summary      Get requests
// @Description  Get requests with filters
// @Tags         requests
// @Accept			 json
// @Produce      json
// @Param 			 request body GetRequestsRequest false "Get requests filters"
// @Success      200 {object} GetRequestsResponse "Requests"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Request not found"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /requests [get] "QUERY"
func (h *HttpHandler) GetRequests(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	var request GetRequestsRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		if !errors.Is(err, core_errors.ErrEmptyRequestBody) {
			responseHandler.ErrorResponse(err, "failed to decode and validate request")
			return
		}
	}

	page, limit, filters := requestsDomainFromDTO(request)

	requests, ceil, err := h.requestsService.GetRequests(ctx, page, limit, filters)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get requests")
		return
	}

	response := requestsDTOFromDomain(requests, ceil)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func requestsDomainFromDTO(request GetRequestsRequest) (int, int, domain.RequestQueryFilter) {
	page, limit := 1, 10

	if request.Page > 0 {
		page = request.Page
	}
	if request.Limit > 0 {
		limit = request.Limit
	}

	filters := domain.RequestQueryFilter{
		Status: request.Statuses,
		Sort:   request.Sort,
	}

	return page, limit, filters
}

func requestsDTOFromDomain(requests []domain.Request, ceil int64) GetRequestsResponse {
	reqs := make([]request_http_dto.GetRequestsResponseRequest, len(requests))
	for i, req := range requests {
		reqs[i] = request_http_dto.RequestsResponseFromDomain(req)
	}

	return GetRequestsResponse{
		Requests:    reqs,
		TotalLength: ceil,
	}
}
