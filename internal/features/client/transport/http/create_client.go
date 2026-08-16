package client_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

type CreateClientRequest struct {
	Company  string   `json:"company"  validate:"required,min=1,max=255"                     example:"CascadePro"`
	Contacts []string `json:"contacts" validate:"required,min=1,dive,required,min=1,max=255" example:""`
}

// CreateClient godoc
// @Summary      Create client
// @Description  Create client in the system
// @Tags         clients
// @Accept       json
// @Param 			 request body CreateClientRequest true "Create client body request"
// @Success      204 "Successfully created client"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /clients [post]
func (h *HttpHandler) CreateClient(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	var request CreateClientRequest
	if err := core_http_request.DecodeAndValidate(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	client := domain.NewCreateClient(request.Company, request.Contacts)

	if _, err := h.clientService.CreateClient(ctx, client); err != nil {
		responseHandler.ErrorResponse(err, "failed to create client")
		return
	}

	responseHandler.NoContentResponse()
}
