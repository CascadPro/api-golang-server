package users_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/cache"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
	"github.com/google/uuid"
)

type Service struct {
	mediaService       media_service.ServiceMethods
	cache              core_redis_cache.Cache
	usersPostgresRepo  UsersRepository
	outboxPostgresRepo OutboxRepository
}

type UsersRepository interface {
	GetUser(ctx context.Context, user domain.User) (domain.User, error)

	PatchUserTx(ctx context.Context, tx core_postgres_pool.Tx, id uuid.UUID, user domain.User) (domain.User, error)

	BeginTx(ctx context.Context) (core_postgres_pool.Tx, error)
}

type OutboxRepository interface {
	CreateEventTx(ctx context.Context, tx core_postgres_pool.Tx, event domain.OutboxEvent) (domain.OutboxEvent, error)
}

type ServiceMethods interface {
	GetCurrentUser(context.Context, uuid.UUID) (domain.User, []byte, error)
	UpdateAvatar(context.Context, uuid.UUID, *domain.File, []byte) error
	DeleteAvatar(context.Context, uuid.UUID) error
}

func NewService(
	mediaService media_service.ServiceMethods,
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
