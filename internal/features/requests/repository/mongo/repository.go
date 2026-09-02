package requests_mongo_repository

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/mongo/pool"
	"github.com/google/uuid"
)

type Repository struct {
	pool core_mongo_pool.Pool
}

type RepositoryMethods interface {
	CreateRequest(ctx context.Context, request domain.Request) (uuid.UUID, error)
	GetRequest(ctx context.Context, id uuid.UUID) (domain.Request, error)
	GetRequests(ctx context.Context, page, limit int, filter domain.RequestQueryFilter) ([]domain.Request, int64, error)
	PatchRequest(ctx context.Context, id uuid.UUID, request domain.Request) error
	PatchRequestStatus(ctx context.Context, id uuid.UUID, userID uuid.UUID, status domain.RequestStatus) error
	DeleteRequest(ctx context.Context, id uuid.UUID) error
}

func NewRepository(pool core_mongo_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
