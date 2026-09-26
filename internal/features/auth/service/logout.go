package auth_service

import (
	"context"
	"fmt"

	core_outbox_utils "github.com/CascadePro/api-golang-server/internal/core/utils/outbox"
)

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	refreshClaims, err := s.tokenIssuer.ParseRefresh(refreshToken)
	if err != nil {
		return fmt.Errorf("parse refresh token: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteSession(ctx, refreshClaims.UserID, refreshClaims.SessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	event, err := core_outbox_utils.NewRealtimeEventSessionRevoked(refreshClaims.UserID, refreshClaims.SessionID)
	if err != nil {
		return err
	}

	if _, err := s.outboxPostgresRepo.CreateEvent(ctx, event); err != nil {
		return fmt.Errorf("create outbox event: %w", err)
	}

	return nil
}
