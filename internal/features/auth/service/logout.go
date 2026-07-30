package auth_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
)

func (s *Service) Logout(ctx context.Context) error {
	userID, err := core_context.UserIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get userID from context: %w", err)
	}

	sessionID, err := core_context.SessionIDFromContext(ctx)
	if err != nil {
		return fmt.Errorf("get sessionID from context: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteSession(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
