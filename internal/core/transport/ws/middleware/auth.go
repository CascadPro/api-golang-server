package core_ws_middleware

import (
	"context"
	"fmt"

	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	sessions_redis_repository "github.com/CascadePro/api-golang-server/internal/features/sessions/repository/redis"
)

type Authenticator struct {
	repository  sessions_redis_repository.RepositoryMethods
	tokenIssuer core_jwt_security.AccessTokenVerifier
}

type AuthenticatorMethods interface {
	Authenticate(context.Context, string) (*core_jwt_security.AccessClaims, error)
}

func NewAuthenticator(repository sessions_redis_repository.RepositoryMethods, tokenIssuer core_jwt_security.AccessTokenVerifier) *Authenticator {
	return &Authenticator{
		repository:  repository,
		tokenIssuer: tokenIssuer,
	}
}

func (a *Authenticator) Authenticate(ctx context.Context, accessToken string) (*core_jwt_security.AccessClaims, error) {
	claims, err := a.tokenIssuer.ParseAccess(accessToken)
	if err != nil {
		return nil, fmt.Errorf("validate websocket session: %w", err)
	}

	if _, err := a.repository.GetSession(ctx, claims.UserID, claims.SessionID); err != nil {
		return nil, fmt.Errorf("validate websocket session: %w", err)
	}

	return claims, nil
}
