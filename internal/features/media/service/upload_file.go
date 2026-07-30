package media_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_media_utils "github.com/CascadePro/api-golang-server/internal/core/utils/media"
)

func (s *Service) UploadFile(ctx context.Context, dtoFile *domain.File, content []byte) (domain.File, error) {
	if err := dtoFile.Validate(); err != nil {
		return domain.File{}, fmt.Errorf("validate file: %w", err)
	}

	file, err := s.mediaPostgresRepo.CreateFile(ctx, dtoFile)
	if err != nil {
		return domain.File{}, fmt.Errorf("create file at repository: %w", err)
	}

	if file.Tag == domain.FileTagAvatars {
		content, err = core_media_utils.ResizeAny(
			content,
			domain.FileAvatarS3Size,
			domain.FileAvatarS3Size,
			&domain.FileAvatarS3Quality,
		)
		if err != nil {
			return domain.File{}, fmt.Errorf("resize image: %w", err)
		}
	}

	key := fmt.Sprintf("%s/%s", file.Tag, file.ID)
	if err := s.coreS3Repo.PutObject(ctx, key, content); err != nil {
		return domain.File{}, fmt.Errorf("put object to S3 repository: %w", err)
	}

	return file, nil
}
