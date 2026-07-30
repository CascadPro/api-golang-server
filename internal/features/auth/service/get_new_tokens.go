package auth_service

import (
	"context"
	"fmt"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
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

	role, err := core_context.UserRoleFromContext(ctx)
	if err != nil {
		return "", fmt.Errorf("get role from context: %w", err)
	}

	accessToken, _, err := core_utils.IssueTokens(claims.UserID, claims.SessionID, role)
	if err != nil {
		return "", fmt.Errorf("issue tokens: %w", err)
	}

	accessTokenString, err := core_utils.GenerateJwtToken(jwt.SigningMethodHS256, accessToken)
	if err != nil {
		return "", fmt.Errorf("generate jwt token string: %w", err)
	}

	return accessTokenString, nil
}
