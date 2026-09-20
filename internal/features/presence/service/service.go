package presence_service

import (
	"context"

	presence_redis_repository "github.com/CascadePro/api-golang-server/internal/features/presence/repository/redis"
	sessions_redis_repository "github.com/CascadePro/api-golang-server/internal/features/sessions/repository/redis"
	"github.com/google/uuid"
)

type Service struct {
	presenceRedisRepo presence_redis_repository.RepositoryMethods
	sessionsRedisRepo sessions_redis_repository.RepositoryMethods
}

type ServiceMethods interface {
	Register(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (RegisterResult, error)
	Unregister(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (UnregisterResult, error)

	GetOnlineSessions(ctx context.Context, userID uuid.UUID, sessionIDs []string) (map[string]bool, error)

	Heartbeat(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) error
}

func NewService(
	presenceRedisRepo presence_redis_repository.RepositoryMethods,
	sessionsRedisRepo sessions_redis_repository.RepositoryMethods,
) *Service {
	return &Service{
		presenceRedisRepo: presenceRedisRepo,
		sessionsRedisRepo: sessionsRedisRepo,
	}
}
