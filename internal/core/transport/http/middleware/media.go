package core_http_middleware

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
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

	allowedTags = []domain.FileTag{domain.FileTagAvatars, domain.FileTagDocs, domain.FileTagImages, domain.FileTagVideos}

	allowedMimesByTag = map[domain.FileTag][]domain.FileMimeType{
		domain.FileTagAvatars: allowedMimes[:4],
		domain.FileTagImages:  allowedMimes[:3],
		domain.FileTagVideos:  allowedMimes[4:5],
		domain.FileTagDocs:    allowedMimes[5:],
	}
)

// MediaMiddleware validates multipart uploads.
//   - The request must have a `Content‑Type: multipart/form-data` header.
//   - A form field named `tag` **must** be present and its value has to be one
//     of the constants from `domain.FileTag`.
//   - At least one file part must be present. The MIME type of the *first* file
//     part is detected (using the part header if it exists, otherwise by
//     sniffing the first 512 bytes) and must belong to the list defined in
//     `domain.FileMimeType`.
//   - When everything is fine the detected mime‑type (and the tag) are stored
//     into the request context so that downstream handlers can read them
//     without having to parse the multipart again.
func Media() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewResponseHandler(log, rw)

			ct := r.Header.Get("Content-Type")
			if ct == "" || !strings.HasPrefix(ct, "multipart/form-data") {
				responseHandler.ErrorResponse(
					wrapValidationError(fmt.Errorf("get `Content-Type` header")),
					"request must be multipart/form-data",
				)
				return
			}

			const maxMemory = 32 << 20 // 32 MB
			if err := r.ParseMultipartForm(maxMemory); err != nil {
				responseHandler.ErrorResponse(
					wrapValidationError(fmt.Errorf("parse multipart form: %v", err)),
					"failed to parse multipart form",
				)
				return
			}

			tagStr := r.FormValue("tag")
			if tagStr == "" {
				responseHandler.ErrorResponse(
					wrapValidationError(fmt.Errorf("form value `tag` can't be NULL")),
					"failed to get form value `tag`",
				)
				return
			}
			tag := domain.FileTag(tagStr)
			if err := core_validation.ValidateArray(allowedTags, tag); err != nil {
				responseHandler.ErrorResponse(
					wrapValidationError(fmt.Errorf("form value tag=%s: %w", tag, err)),
					"failed to validate form value `tag`",
				)
				return
			}

			var foundFileHeader *multipart.FileHeader
			for _, fileHeaders := range r.MultipartForm.File {
				if len(fileHeaders) > 0 {
					foundFileHeader = fileHeaders[0]
					break
				}
			}
			if foundFileHeader == nil {
				responseHandler.ErrorResponse(
					wrapValidationError(fmt.Errorf("get multipart request: file part can't be NULL")),
					"failed to get file from multipart request",
				)
				return
			}

			mimeStr := strings.TrimSpace(foundFileHeader.Header.Get("Content-Type"))

			if mimeStr == "" || mimeStr == "application/octet-stream" {
				f, err := foundFileHeader.Open()
				if err != nil {
					responseHandler.ErrorResponse(fmt.Errorf("open uploaded file: %w", err), "failed to open uploaded file")
					return
				}

				buf := make([]byte, 512)
				n, err := f.Read(buf)
				_ = f.Close()
				if err != nil && err != io.EOF {
					responseHandler.ErrorResponse(fmt.Errorf("read uploaded file: %w", err), "failed to read uploaded file")
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
					fmt.Errorf("file format '%s': %w", tag, err),
					"failed to validate file format",
				)
				return
			}

			if mimes, _ := allowedMimesByTag[tag]; len(mimes) > 0 {
				if err := core_validation.ValidateArray(mimes, mime); err != nil {
					responseHandler.ErrorResponse(
						fmt.Errorf("file format '%s' doesn't satisfy file tag '%s': %w", mime, tag, err),
						"failed to satisfy file tag and file format",
					)
					return
				}
			}

			ctx = context.WithValue(ctx, core_context.CtxKeyMimeType, mime)
			ctx = context.WithValue(ctx, core_context.CtxKeyTag, tag)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}

func wrapValidationError(err error) error {
	return fmt.Errorf("%v: %w", err, core_errors.ErrInvalidArgument)
}
