package users_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	media_service "github.com/CascadePro/api-golang-server/internal/features/media/service"
	users_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/users/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	mediaService      media_service.ServiceMethods
	usersPostgresRepo users_postgres_repository.RepositoryMethods
}

type ServiceMethods interface {
	GetCurrentUser(context.Context, uuid.UUID) (domain.User, error)
}

func NewService(
	mediaService media_service.ServiceMethods,
	usersPostgresRepo users_postgres_repository.RepositoryMethods,
) *Service {
	return &Service{
		mediaService:      mediaService,
		usersPostgresRepo: usersPostgresRepo,
	}
}
