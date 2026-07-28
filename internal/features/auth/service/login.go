package auth_service

import (
	"context"
	"fmt"
	"net"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_utils "github.com/CascadePro/api-golang-server/internal/core/transport/http/utils"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func (s *Service) Login(
	ctx context.Context,
	user domain.User,
	ip net.IP,
	userAgent *core_http_request.UserAgent,
) (string, string, error) {
	userDomain, err := s.userPostgresRepo.GetUser(ctx, user)
	if err != nil {
		return "", "", fmt.Errorf("get user: %w", err)
	}

	if !userDomain.Activated {
		return "", "", fmt.Errorf("you must activate an account before usage: %w", core_errors.ErrConflict)
	}

	isMatching, err := core_utils.CompareStringAndHash(*user.PasswordHash, *userDomain.PasswordHash)
	if err != nil {
		return "", "", fmt.Errorf("compare password hash: %w", err)
	}
	if !isMatching {
		return "", "", fmt.Errorf("compare password hash: %w", core_errors.ErrUnauthorized)
	}

	metadata, err := core_http_utils.GetSessionMetadata(ctx, s.ipinfoRepo, ip, userAgent)
	if err != nil {
		return "", "", fmt.Errorf("get session metadata: %w", err)
	}

	session := domain.NewAuthSession(ip, metadata, domain.SessionLifetime30d)

	sessionID, err := s.sessionsRedisRepo.CreateSession(ctx, userDomain.ID, session)
	if err != nil {
		return "", "", fmt.Errorf("create session: %w", err)
	}

	accessToken, refreshToken, err := issueTokens(userDomain.ID, sessionID, userDomain.Role)
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

func issueTokens(uid uuid.UUID, sid string, role domain.UserRole) (string, string, error) {
	atClaims, rtClaims, err := core_utils.IssueTokens(uid, sid, role)
	if err != nil {
		return "", "", fmt.Errorf("issue tokens: %w", err)
	}

	var method = jwt.SigningMethodHS256

	accessToken, err := core_utils.GenerateJwtToken(method, atClaims)
	if err != nil {
		return "", "", fmt.Errorf("access token generate: %w", err)
	}

	refreshToken, err := core_utils.GenerateJwtToken(method, rtClaims)
	if err != nil {
		return "", "", fmt.Errorf("refresh token generate: %w", err)
	}

	return accessToken, refreshToken, nil
}
