package presence_redis_repository

import (
	"context"
	"fmt"
	"time"

	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"

	"github.com/google/uuid"
)

const (
	LeaseTTL          = 45 * time.Second
	HeartbeatInterval = 20 * time.Second
)

type Repository struct {
	pool core_redis_pool.Pool
}

type RepositoryMethods interface {
	Register(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (bool, error)
	Unregister(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (UnregisterResult, error)

	GetOnlineSessions(ctx context.Context, userID uuid.UUID, sessionIDs []string) (map[string]bool, error)

	Heartbeat(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) error
}

func NewRepository(pool core_redis_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}

func presenceKey(userID uuid.UUID) string {
	return fmt.Sprintf("%s:users:%s", core_redis_pool.PresenceFolder, userID)
}

func presenceMember(sessionID, connectionID string) string {
	return fmt.Sprintf("%s:members:%s:%s", core_redis_pool.PresenceFolder, sessionID, connectionID)
}
