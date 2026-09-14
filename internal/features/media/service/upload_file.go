package media_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func (s *Service) UploadFile(ctx context.Context, uploadedFile *domain.File, content []byte) (domain.File, error) {
	if err := uploadedFile.Validate(); err != nil {
		return domain.File{}, fmt.Errorf("validate file: %w", err)
	}

	file, err := s.mediaPostgresRepo.CreateFile(ctx, uploadedFile)
	if err != nil {
		return domain.File{}, fmt.Errorf("create file at repository: %w", err)
	}

	key := fmt.Sprintf("%s/%s", file.Tag, file.ID)
	if err := s.coreS3Repo.PutObject(ctx, key, content); err != nil {
		_ = s.mediaPostgresRepo.DeleteFile(ctx, file.ID)
		return domain.File{}, fmt.Errorf("put object to S3 repository: %w", err)
	}

	return file, nil
}
