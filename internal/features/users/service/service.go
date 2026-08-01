package users_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	users_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/users/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	usersPostgresRepo users_postgres_repository.RepositoryMethods
}

type ServiceMethods interface {
	GetCurrentUser(context.Context, uuid.UUID) (domain.User, error)
}

func NewService(usersPostgresRepo users_postgres_repository.RepositoryMethods) *Service {
	return &Service{
		usersPostgresRepo: usersPostgresRepo,
	}
}
