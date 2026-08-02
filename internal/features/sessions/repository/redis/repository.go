package sessions_redis_repository

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
	"github.com/google/uuid"
)

type Repository struct {
	pool core_redis_pool.Pool
}

type RepositoryMethods interface {
	CreateSession(ctx context.Context, userID uuid.UUID, session domain.Session) (string, error)
	GetSession(ctx context.Context, userID uuid.UUID, sessionID string) (domain.Session, error)
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, error)
	DeleteSession(ctx context.Context, userID uuid.UUID, sessionID string) error
	DeleteUserSessions(ctx context.Context, userID uuid.UUID, sessionID string) error
}

func NewRepository(pool core_redis_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
