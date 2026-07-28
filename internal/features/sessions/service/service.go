package session_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	sessions_redis_repository "github.com/CascadePro/api-golang-server/internal/features/sessions/repository/redis"
)

type Service struct {
	sessionsRedisRepo sessions_redis_repository.RepositoryMethods
}

type ServiceMethods interface {
	GetUserSessions(context.Context) ([]domain.Session, error)
	DeleteSession(context.Context, string) error
}

func NewService(sessionsRedisRepo sessions_redis_repository.RepositoryMethods) *Service {
	return &Service{
		sessionsRedisRepo: sessionsRedisRepo,
	}
}
