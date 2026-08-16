package client_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	client_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/client/repository/postgres"
	"github.com/google/uuid"
)

type Service struct {
	clientPostgresRepo client_postgres_repository.RepositoryMethods
}

type ServiceMethods interface {
	// ..,
}

func NewService(clientPostgresRepo client_postgres_repository.RepositoryMethods) *Service {
	return &Service{
		clientPostgresRepo: clientPostgresRepo,
	}
}
