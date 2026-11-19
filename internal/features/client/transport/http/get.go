package client_transport_http

import (
	"net/http"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	"github.com/google/uuid"
)

type ClientContactResponse struct {
	Email *string `json:"email,omitempty" example:"test@example.com"`
	Phone *string `json:"phone,omitempty" example:"+79001234567"`
}

type GetClientResponse struct {
	ID        uuid.UUID              `json:"id"         example:"00000000-0000-0000-0000-000000000000"`
	Company   string                 `json:"company"    example:"CascadePro"`
	Contacts  []ClientContactResponse `json:"contacts"`
	CreatedAt time.Time              `json:"created_at" example:"2006-01-02T15:04:05.000000Z"`
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
	locale := core_context.Locale(ctx)

	responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

	clientID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get `id` from request")
		return
	}

	client, err := h.clientService.Get(ctx, clientID)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get client")
		return
	}

	responseHandler.JsonResponse(clientResponseFromDomain(client), http.StatusOK)
}

func clientResponseFromDomain(client domain.Client) GetClientResponse {
	contacts := make([]ClientContactResponse, len(client.Contacts))
	for i, c := range client.Contacts {
		var cc ClientContactResponse
		if c.Email != nil {
			s := c.Email.String()
			cc.Email = &s
		}
		if c.Phone != nil {
			s := c.Phone.String()
			cc.Phone = &s
		}
		contacts[i] = cc
	}

	return GetClientResponse{
		ID:        client.ID,
		Company:   client.Company,
		Contacts:  contacts,
		CreatedAt: client.CreatedAt,
	}
}