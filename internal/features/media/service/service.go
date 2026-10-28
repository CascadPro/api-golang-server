package media_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/s3/pool"
)

type Service struct {
	mediaPostgresRepo  MediaPostgresRepository
	outboxPostgresRepo OutboxRepository
	coreS3Repo         core_s3_pool.Pool
}

type MediaPostgresRepository interface {
	CreateFile(ctx context.Context, file *domain.File) (domain.File, error)

	GetFile(ctx context.Context, fileID string) (domain.File, []byte, error)

	DeleteFile(ctx context.Context, fileID string) error

	DeleteFileTx(ctx context.Context, tx core_postgres_pool.Tx, fileID string) error

	PatchFileTx(ctx context.Context, tx core_postgres_pool.Tx, fileID string, file domain.File) (domain.File, error)

	BeginTx(ctx context.Context) (core_postgres_pool.Tx, error)
}

type OutboxRepository interface {
	CreateEventTx(ctx context.Context, tx core_postgres_pool.Tx, event domain.OutboxEvent) (domain.OutboxEvent, error)
}

type MediaService interface {
	// GetFile(ctx context.Context, fileID string) (file domain.File, err error)
	// GetFileContent(ctx context.Context, fileID string, w, h, quality *int) (
	// 	file domain.File, content []byte, err error,
	// )
	// GetFileStream(ctx context.Context, fileID, byteRange string) (
	// 	file domain.File, body io.ReadCloser, size int64, contentRange string, err error,
	// )
	// GetFilePlaceholder(ctx context.Context, fileID string) (placeholder []byte, err error)

	// UploadFile(ctx context.Context, file *domain.File, content []byte) (created domain.File, err error)

	// DeleteFile(ctx context.Context, tag domain.FileTag, fileID string) error
	// DeleteFileTx(ctx context.Context, tx core_postgres_pool.Tx, tag domain.FileTag, fileID string) error

	// MarkFileDeleted(ctx context.Context, fileID string) error
}

func NewService(
	mediaPostgresRepo MediaPostgresRepository,
	outboxPostgresRepo OutboxRepository,
	coreS3Repo core_s3_pool.Pool,
) *Service {
	return &Service{
		mediaPostgresRepo:  mediaPostgresRepo,
		outboxPostgresRepo: outboxPostgresRepo,
		coreS3Repo:         coreS3Repo,
	}
}
