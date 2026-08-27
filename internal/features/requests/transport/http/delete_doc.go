package requests_transport_http

import (
	"net/http"

	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	request_http_dto "github.com/CascadePro/api-golang-server/internal/features/requests/transport/http/dto"
)

// DeleteDoc godoc
// @Summary      Delete document
// @Description  Delete document and unpin it from request
// @Tags         requests
// @Param        id path string true "Request ID"
// @Param        index path int true "File Index"
// @Success      204 "Successfully deleted document"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 		 409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /requests/{id}/file/{index} [delete]
func (h *HttpHandler) DeleteDoc(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	requestID, index, err := request_http_dto.GetRequestDocPathValues(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request path values")
		return
	}

	if err := h.requestsService.DeleteDoc(ctx, requestID, index); err != nil {
		responseHandler.ErrorResponse(err, "failed to delete document")
		return
	}

	responseHandler.NoContentResponse()
}
