package media_service

import (
	"context"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (s *Service) MarkFileDeleted(ctx context.Context, fileID string) error {
	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return fmt.Errorf("validate file id: %w", err)
	}

	file, err := s.mediaPostgresRepo.GetFile(ctx, fileID)
	if err != nil {
		return fmt.Errorf("get file from repository: %w", err)
	}

	now := time.Now()
	file.Deleted = true
	file.DeletedAt = &now

	if _, err := s.mediaPostgresRepo.PatchFile(ctx, fileID, file); err != nil {
		return fmt.Errorf("patch file in repository")
	}

	key := fmt.Sprintf("%s/%s", file.Tag, fileID)
	if err := s.coreS3Repo.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("delete object from S3 repository: %w", err)
	}

	return nil
}
