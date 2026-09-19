package session_service

import (
	"context"

	core_redis_realtime "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/realtime"
	presence_service "github.com/CascadePro/api-golang-server/internal/features/presence/service"
	sessions_redis_repository "github.com/CascadePro/api-golang-server/internal/features/sessions/repository/redis"
)

type Service struct {
	sessionsRedisRepo sessions_redis_repository.RepositoryMethods
	presenceService   presence_service.ServiceMethods
	publisher         core_redis_realtime.PublisherMethods
}

type ServiceMethods interface {
	GetUserSessions(context.Context) ([]Session, error)
	DeleteSession(context.Context, string) error
	DeleteUserSessions(context.Context) error
}

func NewService(
	sessionsRedisRepo sessions_redis_repository.RepositoryMethods,
	presenceService presence_service.ServiceMethods,
	publisher core_redis_realtime.PublisherMethods,
) *Service {
	return &Service{
		sessionsRedisRepo: sessionsRedisRepo,
		presenceService:   presenceService,
		publisher:         publisher,
	}
}
