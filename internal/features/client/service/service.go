package client_service

import client_postgres_repository "github.com/CascadePro/api-golang-server/internal/features/client/repository/postgres"

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
