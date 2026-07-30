package media_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_s3_pool "github.com/CascadePro/api-golang-server/internal/core/repository/s3/pool"
	media_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/media/repository/postgres"
)

type Service struct {
	mediaPostgresRepo media_postgres_repository.RepositoryMethods
	coreS3Repo        core_s3_pool.Pool
}

type ServiceMethods interface {
	GetFile(context.Context, domain.FileTag, string, *int, *int, *int) (domain.File, []byte, error)
	UploadFile(context.Context, *domain.File, []byte) (domain.File, error)
	DeleteFile(context.Context, domain.FileTag, string) error
	MarkFileDeleted(context.Context, string) error
}

func NewService(
	mediaPostgresRepo media_postgres_repository.RepositoryMethods,
	coreS3Repo core_s3_pool.Pool,
) *Service {
	return &Service{
		mediaPostgresRepo: mediaPostgresRepo,
		coreS3Repo:        coreS3Repo,
	}
}
