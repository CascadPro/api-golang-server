package session_service

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type Service struct {
	sessionsRedisRepo  SessionsRepository
	presenceService    PresenceService
	outboxPostgresRepo OutboxRepository
}

type SessionsRepository interface {
	GetByUser(ctx context.Context, userID uuid.UUID) ([]domain.Session, error)

	Delete(ctx context.Context, userID uuid.UUID, sessionID string) error

	DeleteByUser(ctx context.Context, userID uuid.UUID, sessionID string) error
}

type PresenceService interface {
	GetOnlineSessions(ctx context.Context, userID uuid.UUID, sessionIDs []string) (map[string]bool, error)
}

type OutboxRepository interface {
	CreateEvent(ctx context.Context, event domain.OutboxEvent) (domain.OutboxEvent, error)
}

func NewService(
	sessionsRedisRepo SessionsRepository,
	presenceService PresenceService,
	outboxPostgresRepo OutboxRepository,
) *Service {
	return &Service{
		sessionsRedisRepo:  sessionsRedisRepo,
		presenceService:    presenceService,
		outboxPostgresRepo: outboxPostgresRepo,
	}
}
