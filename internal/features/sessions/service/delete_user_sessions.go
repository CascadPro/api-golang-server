package session_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
)

func (s *Service) DeleteUserSessions(ctx context.Context) error {
	sessionID, err := core_context.SessionIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get sessionID from context: %w", err)
	}

	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get userID from context: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteUserSessions(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("delete user sessions from repository: %w", err)
	}

	return nil
}
