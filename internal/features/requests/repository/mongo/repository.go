package requests_mongo_repository

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/repository/mongo/pool"
	"github.com/google/uuid"
)

type Repository struct {
	pool core_mongo_pool.Pool
}

type RepositoryMethods interface {
	CreateRequest(ctx context.Context, request domain.Request) (uuid.UUID, error)
	GetRequest(ctx context.Context, id uuid.UUID) (domain.Request, error)
	DeleteRequest(ctx context.Context, id uuid.UUID) error
}

func NewRepository(pool core_mongo_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
