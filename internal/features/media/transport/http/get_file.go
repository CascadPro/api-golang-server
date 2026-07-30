package media_transport_http

import (
	"fmt"
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
)

func (h *HttpHandler) GetFile(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewResponseHandler(log, rw)

	fileTag, fileID, err := getFilePathValues(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get path values from request")
		return
	}

	width, height, quality, err := getFileQueryParams(r)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get query params from request")
		return
	}
	if fileTag == domain.FileTagAvatars && width != nil && height != nil {
		if *width > domain.FileAvatarS3Size || *height > domain.FileAvatarS3Size {
			responseHandler.ErrorResponse(fmt.Errorf(
				"size params must be less than or equal %d: %w",
				domain.FileAvatarS3Size, core_errors.ErrInvalidArgument,
			), "invalid width/height query param")
			return
		}
	}

	file, content, err := h.mediaService.GetFile(ctx, fileTag, fileID, width, height, quality)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to get file")
		return
	}
	if content == nil {
		responseHandler.ErrorResponse(core_errors.ErrNotFound, "failed to get file")
	}

	responseHandler.MediaContentResponse(content, file.MimeType)
}

func getFilePathValues(r *http.Request) (domain.FileTag, string, error) {
	fileTag, err := core_http_request.GetFileTagPathValue(r, "tag")
	if err != nil {
		return domain.FileTagNil, "", fmt.Errorf("get path values: %w", err)
	}

	fileID, _, err := core_http_request.GetFileNamePathValue(r, "filename")
	if err != nil {
		return domain.FileTagNil, "", fmt.Errorf("get path values: %w", err)
	}

	return fileTag, fileID, nil
}

func getFileQueryParams(r *http.Request) (*int, *int, *int, error) {
	const defaultError = "get query params: "

	var minimumLen, maximumLen, max = 0, 8192, 250

	width, err := core_http_request.GetIntQueryParam(r, "w", &minimumLen, &maximumLen)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(defaultError+"%w", err)
	}

	height, err := core_http_request.GetIntQueryParam(r, "h", &minimumLen, &maximumLen)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(defaultError+"%w", err)
	}

	quality, err := core_http_request.GetIntQueryParam(r, "quality", &minimumLen, &max)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(defaultError+"%w", err)
	}

	return width, height, quality, nil
}
