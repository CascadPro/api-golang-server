package media_service

import (
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_avatar_utils "github.com/CascadePro/api-golang-server/internal/core/utils/avatar"
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
)

func (s *Service) UploadFile(ctx context.Context, uploadedFile *domain.File, content []byte) (domain.File, error) {
	if uploadedFile == nil {
		return domain.File{}, fmt.Errorf("uploaded file is nil: %w", media_errors.ErrInvalidFile)
	}
	if len(content) == 0 {
		return domain.File{}, fmt.Errorf("uploaded content is empty: %w", media_errors.ErrInvalidFile)
	}

	// All content-specific processing is centralized here. HTTP handlers only
	// parse multipart data; non-HTTP callers get the same validation pipeline.
	if uploadedFile.Tag == domain.FileTagAvatars {
		avatar, err := core_avatar_utils.ProcessAvatar(bytes.NewReader(content))
		if err != nil {
			return domain.File{}, fmt.Errorf("process avatar: %w", err)
		}
		uploadedFile.Filename = replaceFileExtension(uploadedFile.Filename, avatar.Extension)
		uploadedFile.MimeType = avatar.MimeType
		content = avatar.Data
		uploadedFile.Size = int64(len(content))
	} else {
		uploadedFile.Size = int64(len(content))
	}

	if err := validateUploadTagAndMime(uploadedFile.Tag, uploadedFile.MimeType); err != nil {
		return domain.File{}, err
	}
	if err := uploadedFile.Validate(); err != nil {
		return domain.File{}, fmt.Errorf("validate file: %w", err)
	}

	file, err := s.mediaPostgresRepo.CreateFile(ctx, uploadedFile)
	if err != nil {
		return domain.File{}, fmt.Errorf("create file at repository: %w", err)
	}

	key := fmt.Sprintf("%s/%s", file.Tag, file.ID)
	if err := s.coreS3Repo.PutObject(ctx, key, content); err != nil {
		// Best-effort compensation: do not leave metadata for an object that
		// failed to upload. The repository delete is idempotent at this boundary.
		_ = s.mediaPostgresRepo.DeleteFile(ctx, file.ID)
		return domain.File{}, fmt.Errorf("put object to S3 repository: %w", err)
	}
	return file, nil
}

func validateUploadTagAndMime(tag domain.FileTag, mime domain.FileMimeType) error {
	var allowed []domain.FileMimeType
	switch tag {
	case domain.FileTagAvatars:
		allowed = []domain.FileMimeType{domain.FileMimeTypeJpeg, domain.FileMimeTypePng, domain.FileMimeTypeGif, domain.FileMimeTypeWebp}
	case domain.FileTagImages:
		allowed = []domain.FileMimeType{domain.FileMimeTypeJpeg, domain.FileMimeTypePng, domain.FileMimeTypeWebp}
	case domain.FileTagVideos:
		allowed = []domain.FileMimeType{domain.FileMimeTypeMMp4}
	case domain.FileTagDocs:
		allowed = domain.FileMimeTypesDocs
	default:
		return fmt.Errorf("unsupported file tag %q: %w", tag, media_errors.ErrInvalidFile)
	}
	for _, candidate := range allowed {
		if candidate == mime {
			return nil
		}
	}
	return fmt.Errorf("mime type %q is not allowed for tag %q: %w", mime, tag, media_errors.ErrInvalidFile)
}

func replaceFileExtension(filename, newExtension string) string {
	filename = strings.TrimSpace(filename)
	if filename == "" {
		filename = "file"
	}
	ext := filepath.Ext(filename)
	if ext == "." || (strings.HasPrefix(filename, ".") && strings.Count(filename, ".") == 1) {
		ext = ""
	}
	base := strings.TrimSuffix(filename, ext)
	if base == "" {
		base = "file"
	}
	if newExtension != "" && !strings.HasPrefix(newExtension, ".") {
		newExtension = "." + newExtension
	}
	return base + newExtension
}
