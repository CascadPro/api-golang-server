package core_ws_middleware

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	"github.com/google/uuid"
)

type Authenticator struct {
	repository  Repository
	tokenIssuer core_jwt_security.AccessTokenVerifier
}

type Repository interface {
	Get(ctx context.Context, userID uuid.UUID, sessionID string) (domain.Session, error)
}

type AuthenticatorMethods interface {
	Authenticate(context.Context, string) (*core_jwt_security.AccessClaims, error)
}

func NewAuthenticator(repository Repository, tokenIssuer core_jwt_security.AccessTokenVerifier) *Authenticator {
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

	if _, err := a.repository.Get(ctx, claims.UserID, claims.SessionID); err != nil {
		return nil, fmt.Errorf("validate websocket session: %w", err)
	}

	return claims, nil
}
