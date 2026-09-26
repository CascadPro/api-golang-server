package session_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_outbox_utils "github.com/CascadePro/api-golang-server/internal/core/utils/outbox"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (s *Service) DeleteSession(ctx context.Context, sessionID string) error {
	if err := core_validation.ValidateID(sessionID, domain.SessionIDByteLength); err != nil {
		return fmt.Errorf("validate session id: %w", err)
	}

	userID, err := core_context.UserID(ctx)
	if err != nil {
		return fmt.Errorf("get userID from context: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteSession(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("delete session from repository: %w", err)
	}

	event, err := core_outbox_utils.NewRealtimeEventSessionRevoked(userID, sessionID)
	if err != nil {
		return err
	}

	if _, err := s.outboxPostgresRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("create outbox event: %w", err)
	}

	return nil
}
