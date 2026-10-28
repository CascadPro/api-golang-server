package presence_service

import (
	"context"

	presence_redis_repository "github.com/CascadePro/api-golang-server/internal/features/presence/repository/redis"
	"github.com/google/uuid"
)

type Service struct {
	presenceRedisRepo PresenceRepository
	sessionsRedisRepo SessionsRepository
}

type PresenceRepository interface {
	Register(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (presence_redis_repository.RegisterResult, error)

	Unregister(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (presence_redis_repository.UnregisterResult, error)

	GetOnlineSessions(ctx context.Context, userID uuid.UUID, sessionIDs []string) (map[string]bool, error)

	Heartbeat(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) error
}

type SessionsRepository interface {
	PatchLastActive(ctx context.Context, userID uuid.UUID, sessionID string) error
}

type ServiceMethods interface {
	Register(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (RegisterResult, error)
	Unregister(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) (UnregisterResult, error)

	GetOnlineSessions(ctx context.Context, userID uuid.UUID, sessionIDs []string) (map[string]bool, error)

	Heartbeat(ctx context.Context, userID uuid.UUID, sessionID, connectionID string) error
}

func NewService(presenceRedisRepo PresenceRepository, sessionsRedisRepo SessionsRepository) *Service {
	return &Service{
		presenceRedisRepo: presenceRedisRepo,
		sessionsRedisRepo: sessionsRedisRepo,
	}
}
