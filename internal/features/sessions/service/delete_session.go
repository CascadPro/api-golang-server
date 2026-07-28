package session_service

import (
	"context"
	"fmt"

	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (s *Service) DeleteSession(ctx context.Context, sessionID string) error {
	if err := core_validation.ValidateID(sessionID, 24); err != nil {
		return fmt.Errorf("validate session id: %w", err)
	}

	userID, err := core_utils.ParseUUIDFromContext(ctx, "userID")
	if err != nil {
		return fmt.Errorf("parse UUID from context: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteSession(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("delete session from repository: %w", err)
	}

	return nil
}
