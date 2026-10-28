package media_service

import (
	"context"
	"io"

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

type ServiceMethods interface {
	GetFile(context.Context, string) (domain.File, error)
	GetFileContent(context.Context, string, *int, *int, *int) (domain.File, []byte, error)
	GetFileStream(context.Context, string, string) (domain.File, io.ReadCloser, int64, string, error)
	GetFilePlaceholder(context.Context, string) ([]byte, error)

	UploadFile(context.Context, *domain.File, []byte) (domain.File, error)

	DeleteFile(context.Context, domain.FileTag, string) error
	DeleteFileTx(context.Context, core_postgres_pool.Tx, domain.FileTag, string) error

	MarkFileDeleted(context.Context, string) error
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
