package core_http_request

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_avatar_utils "github.com/CascadePro/api-golang-server/internal/core/utils/avatar"
)

func GetFile(r *http.Request, maxSize int64) (multipart.File, *multipart.FileHeader, domain.FileTag, domain.FileMimeType, error) {
	if err := r.ParseMultipartForm(maxSize); err != nil {
		return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf("parse multipart form: %w", err)
	}

	uploadedFile, fileHeader, err := r.FormFile("file")
	if err != nil {
		return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf("retrieve form file: %w", err)
	}
	defer uploadedFile.Close()

	tag, err := core_context.FileTag(r.Context())
	if err != nil {
		return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf("get tag from context: %w", err)
	}

	mimeType, err := core_context.FileMimeType(r.Context())
	if err != nil {
		return nil, nil, domain.FileTagNil, domain.FileMimeTypeNil, fmt.Errorf("get mime type from context: %w", err)
	}

	return uploadedFile, fileHeader, tag, mimeType, nil
}

func GetAvatar(rw http.ResponseWriter, r *http.Request, maxSize int64) (*domain.File, []byte, error) {
	r.Body = http.MaxBytesReader(rw, r.Body, maxSize)

	uploadedFile, fileHeader, tag, _, err := GetFile(r, maxSize)
	if err != nil {
		return nil, nil, fmt.Errorf("get file: %w", err)
	}

	avatar, err := core_avatar_utils.ProcessAvatar(uploadedFile)
	if err != nil {
		return nil, nil, fmt.Errorf("process avatar: %w", err)
	}

	filename := replaceExtension(fileHeader.Filename, avatar.Extension)

	file := domain.NewUploadFile(
		filename,
		tag,
		avatar.MimeType,
		int64(len(avatar.Data)),
	)

	return &file, avatar.Data, nil
}

func GetDocument(rw http.ResponseWriter, r *http.Request, maxSize int64) (*domain.File, []byte, error) {
	r.Body = http.MaxBytesReader(rw, r.Body, maxSize)

	uploadedFile, fileHeader, tag, mimeType, err := GetFile(r, maxSize)
	if err != nil {
		return nil, nil, fmt.Errorf("get file: %w", err)
	}

	content, err := io.ReadAll(uploadedFile)
	if err != nil {
		return nil, nil, fmt.Errorf("read file: %w", err)
	}

	file := domain.NewUploadFile(fileHeader.Filename, tag, mimeType, int64(len(content)))

	return &file, content, nil
}

func replaceExtension(s string, newExtension string) string {
	filenameSlice := strings.Split(s, ".")
	return strings.Join(filenameSlice[:len(filenameSlice)-1], ".") + newExtension
}
