package session_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_outbox_utils "github.com/CascadePro/api-golang-server/internal/core/utils/outbox"
)

func (s *Service) DeleteUserSessions(ctx context.Context) error {
	sessionID, err := core_context.SessionID(ctx)
	if err != nil {
		return fmt.Errorf("get sessionID from context: %w", err)
	}

	userID, err := core_context.UserID(ctx)
	if err != nil {
		return fmt.Errorf("get userID from context: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteUserSessions(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("delete user sessions from repository: %w", err)
	}

	event, err := core_outbox_utils.NewRealtimeEventSessionsRevoked(userID, sessionID)
	if err != nil {
		return err
	}

	if _, err := s.outboxPostgresRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("create outbox event: %w", err)
	}

	return nil
}
