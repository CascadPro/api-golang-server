package auth_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	auth_errors "github.com/CascadePro/api-golang-server/internal/features/auth/errors"
	users_errors "github.com/CascadePro/api-golang-server/internal/features/users/errors"
)

func (s *Service) GetNewTokens(ctx context.Context, token string) (string, error) {
	claims, err := s.tokenIssuer.ParseRefresh(token)
	if err != nil {
		return "", fmt.Errorf("parse refresh token: %v: %w", err, auth_errors.ErrInvalidRefreshToken)
	}

	user, err := s.userPostgresRepo.GetUser(ctx, domain.User{ID: claims.UserID})
	if err != nil {
		return "", fmt.Errorf("get user from repository: %w", err)
	}
	if !user.Activated {
		return "", fmt.Errorf("user is not activated: %w", users_errors.ErrUserNotActivated)
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
