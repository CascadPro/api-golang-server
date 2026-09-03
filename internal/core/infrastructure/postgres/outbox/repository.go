package core_postgres_outbox

import (
	"context"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	"github.com/google/uuid"
)

type Repository struct {
	pool core_postgres_pool.Pool
}

type RepositoryMethods interface {
	CreateEvent(ctx context.Context, tx core_postgres_pool.Tx, event domain.OutboxEvent) (domain.OutboxEvent, error)

	GetPendingEvents(ctx context.Context, workerID uuid.UUID, limit int, ttl time.Duration) ([]domain.OutboxEvent, error)

	MarkEventProcessed(ctx context.Context, id uuid.UUID, workerID uuid.UUID) error

	MarkEventFailed(ctx context.Context, id uuid.UUID, workerID uuid.UUID, err error) error
}

func NewRepository(pool core_postgres_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
