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

	settings, err := s.settingsPostgresRepo.GetUserSettings(ctx, domain.UserSettings{UserID: userDomain.ID})
	if err != nil {
		return "", "", fmt.Errorf("get user settings: %w", err)
	}

	rawDuration, err := core_utils.ParseDurationExtended(string(settings.SessionExpireTerm))
	if err != nil {
		return "", "", fmt.Errorf("parse settings duration: %w", err)
	}

	metadata, err := core_http_utils.GetSessionMetadata(ctx, s.ipinfoRepo, ip, userAgent)
	if err != nil {
		return "", "", fmt.Errorf("get session metadata: %w", err)
	}

	duration := domain.SessionLifetime(rawDuration)
	session := domain.NewAuthSession(ip, metadata, duration)

	sessionID, err := s.sessionsRedisRepo.CreateSession(ctx, userDomain.ID, session)
	if err != nil {
		return "", "", fmt.Errorf("create session: %w", err)
	}

	accessClaims, err := s.tokenIssuer.IssueAccess(userDomain.ID, sessionID, userDomain.Role)
	if err != nil {
		return "", "", fmt.Errorf("issue access token: %w", err)
	}

	refreshClaims, err := s.tokenIssuer.IssueRefresh(userDomain.ID, sessionID, duration)
	if err != nil {
		return "", "", fmt.Errorf("issue refresh token: %w", err)
	}

	accessToken, err := s.tokenIssuer.SignAccess(accessClaims)
	if err != nil {
		return "", "", fmt.Errorf("sign access token: %w", err)
	}

	refreshToken, err := s.tokenIssuer.SignRefresh(refreshClaims)
	if err != nil {
		return "", "", fmt.Errorf("sign refresh token: %w", err)
	}

	return accessToken, refreshToken, nil
}
