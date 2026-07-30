package core_http_request

import (
	"fmt"
	"io"
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func GetFile(rw http.ResponseWriter, r *http.Request, maxSize int64) (*domain.File, []byte, error) {
	r.Body = http.MaxBytesReader(rw, r.Body, maxSize)

	if err := r.ParseMultipartForm(maxSize); err != nil {
		return nil, nil, fmt.Errorf("parse multipart form: %w", err)
	}

	uploadedFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, nil, fmt.Errorf("retrieve form file: %w", err)
	}
	defer uploadedFile.Close()

	content, err := io.ReadAll(uploadedFile)
	if err != nil {
		return nil, nil, fmt.Errorf("read file: %w", err)
	}

	tag, err := core_context.TagFromContext(r.Context())
	if err != nil {
		return nil, nil, fmt.Errorf("get tag from context: %w", err)
	}

	mimeType, err := core_context.MimeTypeFromContext(r.Context())
	if err != nil {
		return nil, nil, fmt.Errorf("get mime type from context: %w", err)
	}

	file := domain.NewUploadFile(fileHeader.Filename, tag, mimeType, int(len(content)))

	if err := file.GeneratePlaceholder(content); err != nil {
		return nil, nil, fmt.Errorf("generate file placeholder: %w", err)
	}

	return &file, content, nil
}
