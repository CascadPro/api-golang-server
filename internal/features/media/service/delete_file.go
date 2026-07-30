package media_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (s *Service) DeleteFile(ctx context.Context, fileTag domain.FileTag, fileID string) error {
	if err := core_validation.ValidateArray(domain.FileTags, fileTag); err != nil {
		return fmt.Errorf("validate file tag: %w", err)
	}

	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return fmt.Errorf("validate file id: %w", err)
	}

	key := fmt.Sprintf("%s/%s", fileTag, fileID)
	if err := s.coreS3Repo.DeleteObject(ctx, key); err != nil {
		return fmt.Errorf("delete object from S3 repository: %w", err)
	}

	if err := s.mediaPostgresRepo.DeleteFile(ctx, fileID); err != nil {
		return fmt.Errorf("delete file from repository")
	}

	return nil
}
