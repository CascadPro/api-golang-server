package auth_service

import (
	"context"
	"fmt"
)

func (s *Service) Logout(ctx context.Context, refreshToken string) error {
	refreshClaims, err := s.tokenIssuer.ParseRefresh(refreshToken)
	if err != nil {
		return fmt.Errorf("parse refresh token: %w", err)
	}

	if err := s.sessionsRedisRepo.DeleteSession(ctx, refreshClaims.UserID, refreshClaims.SessionID); err != nil {
		return fmt.Errorf("delete session: %w", err)
	}

	return nil
}
