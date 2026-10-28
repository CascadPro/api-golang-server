package session_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	presence_service "github.com/CascadePro/api-golang-server/internal/features/presence/service"
	"github.com/google/uuid"
)

type Service struct {
	sessionsRedisRepo  SessionsRepository
	presenceService    presence_service.ServiceMethods
	outboxPostgresRepo OutboxRepository
}

type SessionsRepository interface {
	GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.Session, error)

	Delete(ctx context.Context, userID uuid.UUID, sessionID string) error

	DeleteByUser(ctx context.Context, userID uuid.UUID, sessionID string) error
}

type OutboxRepository interface {
	CreateEvent(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error)
}

type ServiceMethods interface {
	GetUserSessions(context.Context) ([]Session, error)
	DeleteSession(context.Context, string) error
	DeleteUserSessions(context.Context) error
}

func NewService(
	sessionsRedisRepo SessionsRepository,
	presenceService presence_service.ServiceMethods,
	outboxPostgresRepo OutboxRepository,
) *Service {
	return &Service{
		sessionsRedisRepo:  sessionsRedisRepo,
		presenceService:    presenceService,
		outboxPostgresRepo: outboxPostgresRepo,
	}
}
