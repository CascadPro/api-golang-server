package media_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (s *Service) GetFile(ctx context.Context, fileID string) (domain.File, error) {
	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return domain.File{}, fmt.Errorf("validate file id: %w", err)
	}

	file, _, err := s.mediaPostgresRepo.GetFile(ctx, fileID)
	if err != nil {
		return domain.File{}, fmt.Errorf("get file from repository: %w", err)
	}
	if file.Deleted {
		return domain.File{}, fmt.Errorf(
			"file with name '%s' was deleted on %s (UTC): %w",
			file.Filename, file.DeletedAt.Format("02 January, 2006 at 15:04:05"),
			core_errors.ErrNotFound,
		)
	}

	return file, nil
}
