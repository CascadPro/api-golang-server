package users_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_outbox "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/outbox"
	core_redis_cache "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/cache"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
	users_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/users/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	mediaService       media_service.ServiceMethods
	cache              core_redis_cache.Cache
	usersPostgresRepo  users_postgres_repository.RepositoryMethods
	outboxPostgresRepo core_postgres_outbox.RepositoryMethods
}

type ServiceMethods interface {
	GetCurrentUser(context.Context, uuid.UUID) (domain.User, []byte, error)
	UpdateAvatar(context.Context, uuid.UUID, *domain.File, []byte) error
	DeleteAvatar(context.Context, uuid.UUID) error
}

func NewService(
	mediaService media_service.ServiceMethods,
	cache core_redis_cache.Cache,
	usersPostgresRepo users_postgres_repository.RepositoryMethods,
	outboxPostgresRepo core_postgres_outbox.RepositoryMethods,
) *Service {
	return &Service{
		mediaService:       mediaService,
		cache:              cache,
		usersPostgresRepo:  usersPostgresRepo,
		outboxPostgresRepo: outboxPostgresRepo,
	}
}
