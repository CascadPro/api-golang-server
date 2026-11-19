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
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
)

const (
	// The current upload routes accept avatars up to 10 MiB and request
	// documents up to 1 MiB. Leave room for multipart boundaries and headers.
	// This is a request-body limit, not the multipart in-memory threshold.
	maxMultipartRequestSize int64 = 50 << 20
	multipartMaxMemory      int64 = 128 << 20
)

var maxFileSizeByTag = map[domain.FileTag]int64{
	domain.FileTagAvatars: 10 << 20,
	domain.FileTagImages:  50 << 20,
	domain.FileTagVideos:  100 << 20,
	domain.FileTagDocs:    50 << 20,
}

var (
	allowedMimes = []domain.FileMimeType{
		domain.FileMimeTypeJpeg,
		domain.FileMimeTypePng,
		domain.FileMimeTypeWebp,
		domain.FileMimeTypeGif,
		domain.FileMimeTypeBmp,
		domain.FileMimeTypeTiff,
		domain.FileMimeTypeSvg,
		domain.FileMimeTypeHeic,
		domain.FileMimeTypeHeif,
		domain.FileMimeTypeMMp4,
		domain.FileMimeTypeMWebM,
		domain.FileMimeTypeMOgg,
		domain.FileMimeTypeMQuickTime,
		domain.FileMimeTypeMXMsVideo,
		domain.FileMimeTypeMXMatroska,
		domain.FileMimeTypePdf,
		domain.FileMimeTypeDocx,
		domain.FileMimeTypeDoc,
		domain.FileMimeTypeXlsx,
		domain.FileMimeTypeXls,
		domain.FileMimeTypePptx,
		domain.FileMimeTypePpt,
		domain.FileMimeTypeOdt,
		domain.FileMimeTypeOds,
		domain.FileMimeTypeOdp,
		domain.FileMimeTypeTxt,
		domain.FileMimeTypeCsv,
		domain.FileMimeTypeXml,
		domain.FileMimeTypeZip,
		domain.FileMimeTypeGzip,
		domain.FileMimeTypeRar,
		domain.FileMimeType7z,
		domain.FileMimeTypeTar,
	}

	allowedTags = []domain.FileTag{
		domain.FileTagAvatars,
		domain.FileTagDocs,
		domain.FileTagImages,
		domain.FileTagVideos,
	}

	allowedMimesByTag = map[domain.FileTag][]domain.FileMimeType{
		domain.FileTagAvatars: domain.FileMimeTypesPictures[:4], // JPEG, PNG, WebP, GIF only for avatars
		domain.FileTagImages:  domain.FileMimeTypesPictures,
		domain.FileTagVideos:  domain.FileMimeTypesVideos,
		domain.FileTagDocs:    domain.FileMimeTypesDocs,
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
				if err.Error() == "http: request body too large" {
					responseHandler.ErrorResponse(
						fmt.Errorf("request body too large: %w", media_errors.ErrFileTooLarge),
						"request body too large",
					)
					return
				}

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

			// Check file size limit for this tag
			if maxSize, ok := maxFileSizeByTag[tag]; ok {
				if foundFileHeader.Size > maxSize {
					responseHandler.ErrorResponse(
						fmt.Errorf("file size %d exceeds limit %d: %w", foundFileHeader.Size, maxSize, media_errors.ErrFileTooLarge),
						"file too large",
					)
					return
				}
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
					fmt.Errorf("file format '%s': %w", mime, media_errors.ErrUnsupportedFormat),
					"unsupported file format",
				)
				return
			}

			if mimes, exists := allowedMimesByTag[tag]; exists {
				if err := core_validation.ValidateArray(mimes, mime); err != nil {
					allowed := make([]string, len(mimes))
					for i, m := range mimes {
						allowed[i] = string(m)
					}

					responseHandler.ErrorResponse(
						fmt.Errorf("file format '%s' not allowed for tag '%s', allowed: %s: %w", mime, tag, strings.Join(allowed, ", "), media_errors.ErrUnsupportedFormat),
						fmt.Sprintf("unsupported file format for %s, allowed: %s", tag, strings.Join(allowed, ", ")),
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
