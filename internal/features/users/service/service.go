package users_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/cache"
	"github.com/google/uuid"
)

type Service struct {
	mediaService       MediaService
	cache              core_redis_cache.Cache
	usersPostgresRepo  UsersRepository
	outboxPostgresRepo OutboxRepository
}

type MediaService interface {
	GetFilePlaceholder(ctx context.Context, fileID string) (placeholder []byte, err error)

	UploadFile(ctx context.Context, file *domain.File, content []byte) (created domain.File, err error)

	DeleteFile(ctx context.Context, tag domain.FileTag, fileID string) error

	DeleteFileTx(ctx context.Context, tx core_postgres_pool.Tx, tag domain.FileTag, fileID string) error
}

type UsersRepository interface {
	GetUser(ctx context.Context, user domain.User) (domain.User, error)

	PatchUserTx(ctx context.Context, tx core_postgres_pool.Tx, id uuid.UUID, user domain.User) (domain.User, error)

	BeginTx(ctx context.Context) (core_postgres_pool.Tx, error)
}

type OutboxRepository interface {
	CreateEventTx(ctx context.Context, tx core_postgres_pool.Tx, event domain.OutboxEvent) (domain.OutboxEvent, error)
}

func NewService(
	mediaService MediaService,
	cache core_redis_cache.Cache,
	usersPostgresRepo UsersRepository,
	outboxPostgresRepo OutboxRepository,
) *Service {
	return &Service{
		mediaService:       mediaService,
		cache:              cache,
		usersPostgresRepo:  usersPostgresRepo,
		outboxPostgresRepo: outboxPostgresRepo,
	}
}
