package client_transport_http

import (
	"net/http"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	"github.com/google/uuid"
)

type GetClientResponse struct {
	ID        uuid.UUID `json:"id"         example:"00000000-000000-000000-000000000000"`
	Company   string    `json:"company"    example:"CascadePro"`
	Contacts  []string  `json:"contacts"   example:""`
	CreatedAt time.Time `json:"created_at" example:"2006-01-02T15-04-05.000000"`
}

// GetClient godoc
// @Summary      Get client
// @Description  Get client by ID
// @Tags         clients
// @Param        id path string true "Client ID"
// @Produce      json
// @Success      200 {object} GetClientResponse "Request"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Client not found"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /clients/{id} [get]
func (h *HttpHandler) GetClient(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	clientID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get `id` from request")
		return
	}

	client, err := h.clientService.GetClient(ctx, clientID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get client")
		return
	}

	response := getClientResponseFromDomain(client)

	responseHandler.JsonResponse(response, http.StatusOK)
}

func getClientResponseFromDomain(client domain.Client) GetClientResponse {
	return GetClientResponse{
		ID:        client.ID,
		Company:   client.Company,
		Contacts:  client.Contacts,
		CreatedAt: client.CreatedAt,
	}
}
