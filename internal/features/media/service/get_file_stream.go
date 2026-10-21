package media_service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/s3/pool"
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
)

func (s *Service) GetFileStream(ctx context.Context, fileID, byteRange string) (domain.File, io.ReadCloser, int64, string, error) {
	file, err := s.GetFile(ctx, fileID)
	if err != nil {
		return domain.File{}, nil, 0, "", fmt.Errorf("get file metadata: %w", err)
	}

	key := fmt.Sprintf("%s/%s", file.Tag, file.ID)
	body, size, contentRange, err := s.coreS3Repo.GetObjectStream(ctx, key, byteRange)

	if err != nil {
		if errors.Is(err, core_s3_pool.ErrNotFound) {
			return domain.File{}, nil, 0, "", fmt.Errorf("file %s is missing in S3: %w", fileID, media_errors.ErrFileNotFound)
		}
		return domain.File{}, nil, 0, "", fmt.Errorf("stream file from S3: %w", err)
	}

	// S3's response length is authoritative for ranges; for a full response
	// metadata is the fallback for compatible S3 implementations.
	if byteRange == "" && size == 0 {
		size = file.Size
	}

	return file, body, size, contentRange, nil
}
