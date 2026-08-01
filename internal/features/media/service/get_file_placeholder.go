package media_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (s *Service) GetFilePlaceholder(ctx context.Context, fileID string) ([]byte, error) {
	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return nil, fmt.Errorf("validate file id: %w", err)
	}

	file, placeholder, err := s.mediaPostgresRepo.GetFile(ctx, fileID)
	if err != nil {
		return nil, fmt.Errorf("get file from repository: %w", err)
	}
	if file.Tag != domain.FileTagAvatars && file.Tag != domain.FileTagImages && file.Tag != domain.FileTagVideos {
		return nil, fmt.Errorf("this type of file don't have placeholder: %w", core_errors.ErrInvalidArgument)
	}
	if len(placeholder) <= 0 {
		return nil, fmt.Errorf("placeholder bytes is empty: %w", core_errors.ErrNotFound)
	}

	return placeholder, nil
}
