package requests_transport_http

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	request_http_dto "github.com/CascadePro/api-golang-server/internal/features/requests/transport/http/dto"
)

// UploadDoc godoc
// @Summary      Upload document
// @Description  Upload document and pin it to request
// @Tags         requests
// @Accept       mpfd
// @Param        file formData file true "File to upload"
// @Param        tag formData string true "Text tag"
// @Param        idx path int true "File Index"
// @Success      204 "Successfully uploaded document"
// @Failure      400 {object} core_http_response.ErrorResponse "Bad request"
// @Failure      401 {object} core_http_response.ErrorResponse "Unauthorized"
// @Failure      404 {object} core_http_response.ErrorResponse "Not found"
// @Failure 		 409 {object} core_http_response.ErrorResponse "Conflict error"
// @Failure      429 {object} core_http_response.ErrorResponse "Too many requests"
// @Failure      500 {object} core_http_response.ErrorResponse "Internal server error"
// @Router       /requests/{id}/file/{index} [post]
func (h *HttpHandler) UploadDoc(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	requestID, index, err := request_http_dto.GetRequestDocPathValues(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request path values")
		return
	}

	uploadedFile, content, err := core_http_request.GetFile(rw, r, 1<<20)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get file from request")
		return
	}
	if uploadedFile.Tag != domain.FileTagDocs {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "invalid tag, expected 'docs'")
		return
	}

	if err := h.requestsService.UploadDoc(ctx, requestID, *uploadedFile, content, index); err != nil {
		responseHandler.ErrorResponse(err, "failed to upload doc")
		return
	}

	responseHandler.NoContentResponse()
}
