package core_http_request

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path/filepath"
	"strings"
	"unicode"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

const (
	defaultFilename  = "file"
	defaultExtension = "bin"
)

func GetFile(r *http.Request, maxSize int64) (multipart.File, *multipart.FileHeader, domain.FileTag, domain.FileMimeType, error) {
	// Media middleware parses the multipart body first. Keep this helper usable
	// on routes without that middleware, but do not parse the same form twice.
	if r.MultipartForm == nil {
		if err := r.ParseMultipartForm(32 << 20); err != nil {
			return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf("parse multipart form: %w", err)
		}
	}

	headers := r.MultipartForm.File["file"]
	if len(headers) != 1 {
		return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf(
			"multipart field `file` must contain exactly one file: %w",
			core_errors.ErrInvalidArgument,
		)
	}

	fileHeader := headers[0]
	if maxSize > 0 && fileHeader.Size > maxSize {
		return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf(
			"uploaded file size %d exceeds limit %d: %w",
			fileHeader.Size, maxSize, core_errors.ErrInvalidArgument,
		)
	}

	uploadedFile, err := fileHeader.Open()
	if err != nil {
		return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf("open multipart file: %w", err)
	}

	tag, err := core_context.FileTag(r.Context())
	if err != nil {
		_ = uploadedFile.Close()
		return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf("get tag from context: %w", err)
	}

	mimeType, err := core_context.FileMimeType(r.Context())
	if err != nil {
		_ = uploadedFile.Close()
		return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf("get mime type from context: %w", err)
	}

	return uploadedFile, fileHeader, tag, mimeType, nil
}

// GetAvatar and GetDocument are intentionally transport-only: they read the
// upload and construct metadata. Format conversion and content validation live
// in the media service, shared by every caller of UploadFile.
func GetAvatar(rw http.ResponseWriter, r *http.Request, maxSize int64) (*domain.File, []byte, error) {
	uploadedFile, fileHeader, tag, mimeType, err := GetFile(r, maxSize)
	if err != nil {
		return nil, nil, fmt.Errorf("get file: %w", err)
	}
	defer uploadedFile.Close()

	content, err := io.ReadAll(io.LimitReader(uploadedFile, maxSize+1))
	if err != nil {
		return nil, nil, fmt.Errorf("read avatar upload: %w", err)
	}

	if maxSize > 0 && int64(len(content)) > maxSize {
		return nil, nil, fmt.Errorf("uploaded file exceeds limit %d: %w", maxSize, core_errors.ErrInvalidArgument)
	}

	if len(content) == 0 {
		return nil, nil, fmt.Errorf("uploaded file is empty: %w", core_errors.ErrInvalidArgument)
	}

	file := domain.NewUploadFile(fileHeader.Filename, tag, mimeType, int64(len(content)))

	return &file, content, nil
}

func GetDocument(rw http.ResponseWriter, r *http.Request, maxSize int64) (*domain.File, []byte, error) {
	uploadedFile, fileHeader, tag, mimeType, err := GetFile(r, maxSize)
	if err != nil {
		return nil, nil, fmt.Errorf("get file: %w", err)
	}
	defer uploadedFile.Close()

	content, err := io.ReadAll(io.LimitReader(uploadedFile, maxSize+1))
	if err != nil {
		return nil, nil, fmt.Errorf("read document upload: %w", err)
	}

	if maxSize > 0 && int64(len(content)) > maxSize {
		return nil, nil, fmt.Errorf("uploaded file exceeds limit %d: %w", maxSize, core_errors.ErrInvalidArgument)
	}

	if len(content) == 0 {
		return nil, nil, fmt.Errorf("uploaded file is empty: %w", core_errors.ErrInvalidArgument)
	}

	file := domain.NewUploadFile(fileHeader.Filename, tag, mimeType, int64(len(content)))

	return &file, content, nil
}

func replaceExtension(filename string, newExtension string) string {
	filename = strings.TrimSpace(filename)
	newExtension = strings.TrimSpace(newExtension)

	if filename == "" {
		filename = defaultFilename
	}

	newExtension = strings.TrimLeft(newExtension, ".")

	var ext strings.Builder
	for _, r := range newExtension {
		switch {
		case unicode.IsLetter(r), unicode.IsDigit(r), r == '-', r == '_':
			ext.WriteRune(r)
		default:
			ext.WriteByte('_')
		}
	}

	newExtension = strings.Trim(ext.String(), "_")

	if newExtension == "" {
		newExtension = defaultExtension
	}

	dir := filepath.Dir(filename)
	base := filepath.Base(filename)

	oldExt := filepath.Ext(base)

	if oldExt == base && strings.HasPrefix(base, ".") {
		oldExt = ""
	}

	name := strings.TrimSuffix(base, oldExt)

	if name == "" || name == "." || name == ".." {
		name = "file"
	}

	result := name + "." + newExtension

	if dir == "." {
		return result
	}

	return filepath.Join(dir, result)
}
