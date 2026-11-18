package client_transport_http

import (
	"errors"
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type GetClientsRequest struct {
	Limit  int     `json:"limit" validate:"omitempty,number,gte=1,lte=100" example:"10"`
	Page   int     `json:"page"  validate:"omitempty,number,gte=1,lte=100" example:"1"`
	Search *string `json:"q"     validate:"omitempty,max=100"              example:"search query"`
}

type GetClientsResponse struct {
	Clients []GetClientResponse `json:"items"`
	Total   int64               `json:"length" example:"100"`
}

// GetClients godoc
// @Summary      Get clients
// @Description  Get clients with filters
// @Tags         clients
// @Accept			 json
// @Produce      json
// @Param 			 request body GetClientsRequest false "Get clients filters"
// @Success      200 {object} GetClientsResponse "Clients list"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Request not found"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /clients [get] "QUERY"
func (h *HttpHandler) GetClients(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	locale := core_context.Locale(ctx)

	responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

	var request GetClientsRequest
	if err := core_http_request.DecodeAndValidate(rw, r, &request); err != nil {
		if !errors.Is(err, core_errors.ErrEmptyRequestBody) {
			responseHandler.ErrorResponse(err, "failed to decode and validate request")
			return
		}
	}

	page, limit, searchQuery := clientsDomainFromDTO(request)

	result, err := h.clientService.GetMany(ctx, page, limit, searchQuery)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get clients")
		return
	}

	response := GetClientsResponse{
		Clients: clientsResponseFromDomains(result.Clients),
		Total:   result.Length,
	}

	responseHandler.JsonResponse(response, http.StatusOK)
}

func clientsDomainFromDTO(request GetClientsRequest) (int64, int64, *string) {
	page, limit := 1, 10

	if request.Page > 0 {
		page = request.Page
	}
	if request.Limit > 0 {
		limit = request.Limit
	}

	return int64(page), int64(limit), request.Search
}

func clientsResponseFromDomains(clients []domain.Client) []GetClientResponse {
	response := make([]GetClientResponse, len(clients))
	for i, client := range clients {
		response[i] = clientResponseFromDomain(client)
	}

	return response
}
