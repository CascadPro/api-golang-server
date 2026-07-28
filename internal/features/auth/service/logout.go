package auth_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
)

func (s *Service) Logout(ctx context.Context) error {
	userID, err := core_utils.ParseUUIDFromContext(ctx, "userID")
	if err != nil {
		return fmt.Errorf("parse userID from context: %w", err)
	}

	sessionID, err := core_utils.ParseIDStingFromContext(ctx, "sessionID", domain.SessionIDByteLength)
	if err != nil {
		return fmt.Errorf("parse session id from context: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteSession(ctx, userID, sessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
