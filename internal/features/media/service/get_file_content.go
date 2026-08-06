package media_service

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/repository/s3/pool"
	core_media_utils "github.com/CascadePro/api-golang-server/internal/core/utils/media"
)

func (s *Service) GetFileContent(
	ctx context.Context,
	fileTag domain.FileTag,
	fileID string,
	w, h, quality *int,
) (domain.File, []byte, error) {
	file, err := s.GetFile(ctx, fileTag, fileID)
	if err != nil {
		return domain.File{}, nil, fmt.Errorf("get file from repository: %w", err)
	}

	key := fmt.Sprintf("%s/%s", fileTag, fileID)

	result, err := s.coreS3Repo.GetObject(ctx, key)
	if err != nil {
		if errors.Is(err, core_s3_pool.ErrNotFound) {
			return domain.File{}, nil, fmt.Errorf("file with id=%s in S3 bucket: %w", fileID, core_errors.ErrNotFound)
		}

		return domain.File{}, nil, err
	}

	if w != nil && h != nil && (fileTag == domain.FileTagAvatars || fileTag == domain.FileTagImages) {
		result, err = core_media_utils.ResizeAny(result, *w, *h, quality)
		if err != nil {
			return domain.File{}, nil, fmt.Errorf("resize image file: %w", err)
		}
	}

	return file, result, nil
}
