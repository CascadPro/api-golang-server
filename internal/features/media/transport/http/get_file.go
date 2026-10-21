package media_transport_http

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/s3/pool"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_http_utils "github.com/CascadePro/api-golang-server/internal/core/transport/http/utils"
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
	"go.uber.org/zap"
)

func (h *HttpHandler) GetFile(rw http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	locale := core_context.Locale(ctx)

	responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

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
			responseHandler.ErrorResponse(fmt.Errorf("size params must be less than or equal %d: %w", domain.FileAvatarS3Size, core_errors.ErrInvalidArgument), "invalid width/height query param")

			return
		}
	}

	// A resized image is a transformed representation. Serve it from memory;
	// stream original objects and videos directly from S3, including byte ranges.
	transform := width != nil && height != nil && (fileTag == domain.FileTagAvatars || fileTag == domain.FileTagImages)
	if transform {
		file, content, err := h.mediaService.GetFileContent(ctx, fileID, width, height, quality)
		if err != nil {
			responseHandler.ErrorResponse(err, "failed to get transformed file")
			return
		}

		if file.Tag != fileTag {
			responseHandler.ErrorResponse(media_errors.ErrFileNotFound, "file tag does not match URL")
			return
		}

		setFileHeaders(rw, file)

		rw.Header().Set("Content-Length", strconv.Itoa(len(content)))
		rw.WriteHeader(http.StatusOK)

		_, _ = rw.Write(content)

		return
	}

	rangeHeader := strings.TrimSpace(r.Header.Get("Range"))

	file, body, size, contentRange, err := h.mediaService.GetFileStream(ctx, fileID, rangeHeader)
	if err != nil {
		if errors.Is(err, core_s3_pool.ErrInvalidRange) {
			rw.Header().Set("Content-Range", fmt.Sprintf("bytes */%d", size))
			rw.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
			return
		}

		responseHandler.ErrorResponse(err, "failed to stream file")
		return
	}
	defer body.Close()

	if file.Tag != fileTag {
		responseHandler.ErrorResponse(media_errors.ErrFileNotFound, "file tag does not match URL")
		return
	}

	setFileHeaders(rw, file)

	rw.Header().Set("Accept-Ranges", "bytes")
	rw.Header().Set("Content-Length", strconv.FormatInt(size, 10))

	status := http.StatusOK

	if rangeHeader != "" && contentRange != "" {
		rw.Header().Set("Content-Range", contentRange)
		status = http.StatusPartialContent
	}

	rw.WriteHeader(status)

	if r.Method != http.MethodHead {
		if _, err := io.Copy(rw, body); err != nil {
			log.Warn("streaming media response failed", zap.Error(err))
		}
	}
}

func setFileHeaders(rw http.ResponseWriter, file domain.File) {
	rw.Header().Set("Content-Type", string(file.MimeType))
	rw.Header().Set("X-Content-Type-Options", "nosniff")
	rw.Header().Set("ETag", fmt.Sprintf(`"%s-%d"`, file.ID, file.Version))

	if file.Tag == domain.FileTagAvatars {
		core_http_utils.SetCacheControl(rw, core_http_utils.CacheControlOptions{
			Type:      core_http_utils.CacheControlTypePublic,
			MaxAge:    24 * time.Hour,
			Immutable: true,
		})
	} else {
		rw.Header().Set("Cache-Control", "private, no-store")
	}
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

	minLen, maxLen, maxValue := 0, 8192, 250

	width, err := core_http_request.GetIntQueryParam(r, "w", &minLen, &maxLen)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(defaultError+"%w", err)
	}

	height, err := core_http_request.GetIntQueryParam(r, "h", &minLen, &maxLen)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(defaultError+"%w", err)
	}

	quality, err := core_http_request.GetIntQueryParam(r, "quality", &minLen, &maxValue)
	if err != nil {
		return nil, nil, nil, fmt.Errorf(defaultError+"%w", err)
	}

	return width, height, quality, nil
}
