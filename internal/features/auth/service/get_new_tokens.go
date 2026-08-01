package auth_service

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
	"github.com/golang-jwt/jwt/v5"
)

func (s *Service) GetNewTokens(ctx context.Context, token string) (string, error) {
	claims, err := core_utils.GetJwtTokenClaims[core_utils.JwtRefreshClaims](token)
	if err != nil {
		return "", fmt.Errorf("get jwt token claims: %w", err)
	}

	_, err = s.sessionsRedisRepo.GetSession(ctx, claims.UserID, claims.SessionID)
	if err != nil {
		return "", fmt.Errorf("get session: %w", err)
	}

	user, err := s.userPostgresRepo.GetUser(ctx, domain.User{ID: claims.UserID})
	if err != nil {
		return "", fmt.Errorf("get user from repository: %w", err)
	}
	if !user.Activated {
		return "", fmt.Errorf("user is not activated: %w", core_errors.ErrConflict)
	}

	accessToken, _, err := core_utils.IssueTokens(claims.UserID, claims.SessionID, user.Role)
	if err != nil {
		return "", fmt.Errorf("issue tokens: %w", err)
	}

	accessTokenString, err := core_utils.GenerateJwtToken(jwt.SigningMethodHS256, accessToken)
	if err != nil {
		return "", fmt.Errorf("generate jwt token string: %w", err)
	}

	return accessTokenString, nil
}
