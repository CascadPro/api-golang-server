package client_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type Service struct {
	clientPostgresRepo ClientPostgresRepository
}

type ClientPostgresRepository interface {
	Create(ctx context.Context, client domain.Client) (domain.Client, error)
	Get(ctx context.Context, id uuid.UUID) (domain.Client, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type ServiceMethods interface {
	GetClient(context.Context, uuid.UUID) (domain.Client, error)
	CreateClient(context.Context, domain.Client) (domain.Client, error)
	DeleteClient(context.Context, uuid.UUID) error
}

func NewService(clientPostgresRepo ClientPostgresRepository) *Service {
	return &Service{
		clientPostgresRepo: clientPostgresRepo,
	}
}
