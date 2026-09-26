package session_service

import (
	"context"

	core_postgres_outbox "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/outbox"
	presence_service "github.com/CascadePro/api-golang-server/internal/features/presence/service"
	sessions_redis_repository "github.com/CascadePro/api-golang-server/internal/features/sessions/repository/redis"
)

type Service struct {
	sessionsRedisRepo  sessions_redis_repository.RepositoryMethods
	presenceService    presence_service.ServiceMethods
	outboxPostgresRepo core_postgres_outbox.RepositoryMethods
}

type ServiceMethods interface {
	GetUserSessions(context.Context) ([]Session, error)
	DeleteSession(context.Context, string) error
	DeleteUserSessions(context.Context) error
}

func NewService(
	sessionsRedisRepo sessions_redis_repository.RepositoryMethods,
	presenceService presence_service.ServiceMethods,
	outboxPostgresRepo core_postgres_outbox.RepositoryMethods,
) *Service {
	return &Service{
		sessionsRedisRepo:  sessionsRedisRepo,
		presenceService:    presenceService,
		outboxPostgresRepo: outboxPostgresRepo,
	}
}
