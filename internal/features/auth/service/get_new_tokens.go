package auth_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

func (s *Service) GetNewTokens(ctx context.Context, token string) (string, error) {
	claims, err := s.tokenIssuer.ParseRefresh(token)
	if err != nil {
		return "", fmt.Errorf("parse refresh token: %w", err)
	}

	user, err := s.userPostgresRepo.GetUser(ctx, domain.User{ID: claims.UserID})
	if err != nil {
		return "", fmt.Errorf("get user from repository: %w", err)
	}
	if !user.Activated {
		return "", fmt.Errorf("user is not activated: %w", core_errors.ErrConflict)
	}

	accessClaims, err := s.tokenIssuer.IssueAccess(claims.UserID, claims.SessionID, user.Role)
	if err != nil {
		return "", fmt.Errorf("issue access token: %w", err)
	}

	accessToken, err := s.tokenIssuer.SignAccess(accessClaims)
	if err != nil {
		return "", fmt.Errorf("sign access token: %w", err)
	}

	return accessToken, nil
}
