package core_http_middleware

import (
	"fmt"
	"io"
	"mime"
	"net/http"
	"strings"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

const (
	// The current upload routes accept avatars up to 10 MiB and request
	// documents up to 1 MiB. Leave room for multipart boundaries and headers.
	// This is a request-body limit, not the multipart in-memory threshold.
	maxMultipartRequestSize int64 = 12 << 20
	multipartMaxMemory      int64 = 32 << 20
)

var (
	allowedMimes = []domain.FileMimeType{
		domain.FileMimeTypeJpeg,
		domain.FileMimeTypePng,
		domain.FileMimeTypeWebp,
		domain.FileMimeTypeGif,
		domain.FileMimeTypeMMp4,
		domain.FileMimeTypePdf,
		domain.FileMimeTypeDocx,
		domain.FileMimeTypeXlsx,
		domain.FileMimeTypePptx,
	}

	allowedTags = []domain.FileTag{
		domain.FileTagAvatars,
		domain.FileTagDocs,
		domain.FileTagImages,
		domain.FileTagVideos,
	}

	allowedMimesByTag = map[domain.FileTag][]domain.FileMimeType{
		domain.FileTagAvatars: allowedMimes[:4],
		domain.FileTagImages:  allowedMimes[:3],
		domain.FileTagVideos:  allowedMimes[4:5],
		domain.FileTagDocs:    allowedMimes[5:],
	}
)

// Media validates multipart uploads and stores the validated tag and MIME type
// in the request context. Exactly one uploaded file is accepted, under the
// field name "file".
func Media() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			locale := core_context.Locale(ctx)
			responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

			ct := r.Header.Get("Content-Type")
			mediaType, _, parseErr := mime.ParseMediaType(ct)
			if parseErr != nil || !strings.EqualFold(mediaType, "multipart/form-data") {
				responseHandler.ErrorResponse(
					fmt.Errorf("invalid `Content-Type` header: %w", core_errors.ErrInvalidArgument),
					"request must be multipart/form-data",
				)
				return
			}

			// The body limit must be installed before ParseMultipartForm reads
			// the request. ParseMultipartForm's maxMemory is not a body limit.
			r.Body = http.MaxBytesReader(rw, r.Body, maxMultipartRequestSize)
			if err := r.ParseMultipartForm(multipartMaxMemory); err != nil {
				responseHandler.ErrorResponse(
					fmt.Errorf("parse multipart form: %w", err),
					"failed to parse multipart form",
				)
				return
			}
			if r.MultipartForm != nil {
				defer r.MultipartForm.RemoveAll()
			}

			tagStr := r.FormValue("tag")
			if tagStr == "" {
				responseHandler.ErrorResponse(
					fmt.Errorf("form value `tag` can't be NULL: %w", core_errors.ErrInvalidArgument),
					"failed to get form value `tag`",
				)
				return
			}

			tag := domain.FileTag(tagStr)
			if err := core_validation.ValidateArray(allowedTags, tag); err != nil {
				responseHandler.ErrorResponse(
					fmt.Errorf("form value tag=%s: %w", tag, err),
					"failed to validate form value `tag`",
				)
				return
			}

			fileHeaders, ok := r.MultipartForm.File["file"]
			if !ok || len(fileHeaders) != 1 {
				responseHandler.ErrorResponse(
					fmt.Errorf("multipart request must contain exactly one file in field `file`: %w", core_errors.ErrInvalidArgument),
					"failed to get file from multipart request",
				)
				return
			}

			// Reject additional file fields. Otherwise the middleware could
			// validate one file while the handler consumes another.
			for field, headers := range r.MultipartForm.File {
				if field != "file" || len(headers) != 1 {
					responseHandler.ErrorResponse(
						fmt.Errorf("unexpected multipart file field %q: %w", field, core_errors.ErrInvalidArgument),
						"invalid multipart file fields",
					)
					return
				}
			}

			foundFileHeader := fileHeaders[0]
			if foundFileHeader == nil {
				responseHandler.ErrorResponse(
					fmt.Errorf("file part can't be NULL: %w", core_errors.ErrInvalidArgument),
					"failed to get file from multipart request",
				)
				return
			}

			mimeStr := strings.TrimSpace(foundFileHeader.Header.Get("Content-Type"))
			if mimeStr == "" || strings.EqualFold(mimeStr, "application/octet-stream") {
				f, err := foundFileHeader.Open()
				if err != nil {
					responseHandler.ErrorResponse(fmt.Errorf("open uploaded file: %w", err), "failed to open uploaded file")
					return
				}

				buf := make([]byte, 512)
				n, readErr := f.Read(buf)
				closeErr := f.Close()

				if readErr != nil && readErr != io.EOF {
					responseHandler.ErrorResponse(fmt.Errorf("read uploaded file: %w", readErr), "failed to read uploaded file")
					return
				}

				if closeErr != nil {
					responseHandler.ErrorResponse(fmt.Errorf("close uploaded file: %w", closeErr), "failed to close uploaded file")
					return
				}

				mimeStr = http.DetectContentType(buf[:n])
			}

			if idx := strings.Index(mimeStr, ";"); idx >= 0 {
				mimeStr = strings.TrimSpace(mimeStr[:idx])
			}

			mime := domain.FileMimeType(mimeStr)

			if err := core_validation.ValidateArray(allowedMimes, mime); err != nil {
				responseHandler.ErrorResponse(
					fmt.Errorf("file format '%s': %w", mime, err),
					"failed to validate file format",
				)
				return
			}

			if mimes, exists := allowedMimesByTag[tag]; exists {
				if err := core_validation.ValidateArray(mimes, mime); err != nil {
					responseHandler.ErrorResponse(
						fmt.Errorf("file format '%s' doesn't satisfy file tag '%s': %w", mime, tag, err),
						"failed to satisfy file tag and file format",
					)
					return
				}
			}

			ctx = core_context.WithFileMimeType(ctx, mime)
			ctx = core_context.WithFileTag(ctx, tag)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
