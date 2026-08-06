package requests_transport_http

import (
	"fmt"
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

func (h *HttpHandler) UploadDoc(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	log.Info("1")

	requestID, err := core_http_request.GetUUIDPathValue(r, "id")
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get request id path value")
		return
	}

	log.Info("2")

	index, err := getQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get `idx` query param")
		return
	}

	log.Info("3")

	uploadedFile, content, err := core_http_request.GetFile(rw, r, 1<<20)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get file from request")
		return
	}
	if uploadedFile.Tag != domain.FileTagDocs {
		responseHandler.ErrorResponse(core_errors.ErrInvalidArgument, "invalid tag, expected 'docs'")
		return
	}

	log.Info("4")

	if err := h.requestsService.UploadDoc(ctx, requestID, *uploadedFile, content, index); err != nil {
		responseHandler.ErrorResponse(err, "failed to upload doc")
		return
	}

	log.Info("5")

	responseHandler.NoContentResponse()
}

var docNames = []string{"project", "tech_task", "specification", "contract"}

func getQueryParams(r *http.Request) (int, error) {
	var min, max = 0, 3

	index, err := core_http_request.GetIntQueryParam(r, "idx", &min, &max)
	if err != nil {
		return -1, fmt.Errorf("get int query param: %w", err)
	}
	if index == nil {
		return -1, fmt.Errorf("`idx` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	return *index, nil
}
