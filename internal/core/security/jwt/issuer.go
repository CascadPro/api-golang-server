package core_jwt_security

import (
	"fmt"
	"time"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	AccessTokenLifetime time.Duration = 15 * time.Minute

	AuthTokenIssuer = "cascade-pro-api-server"
)

type Issuer struct {
	secret        []byte
	signingMethod jwt.SigningMethod
}

type IssuerMethods interface {
	IssueAccess(userID uuid.UUID, sessionID string, role domain.UserRole) (*AccessClaims, error)
	IssueRefresh(userID uuid.UUID, sessionID string, lifetime domain.SessionLifetime) (*RefreshClaims, error)

	SignAccess(claims *AccessClaims) (string, error)
	SignRefresh(claims *RefreshClaims) (string, error)

	ParseAccess(tokenString string) (*AccessClaims, error)
	ParseRefresh(tokenString string) (*RefreshClaims, error)
}

type AccessTokenVerifier interface {
	ParseAccess(tokenString string) (*AccessClaims, error)
}

func NewIssuer(method jwt.SigningMethod) (*Issuer, error) {
	cfg, err := core_config.NewSecretConfig()
	if err != nil {
		return nil, fmt.Errorf("get jwt secret config: %w", err)
	}

	return &Issuer{
		secret:        []byte(cfg.JwtSecretKey),
		signingMethod: method,
	}, nil
}

func (i *Issuer) IssueAccess(userID uuid.UUID, sessionID string, role domain.UserRole) (*AccessClaims, error) {
	now := time.Now()

	claims := &AccessClaims{
		UserID:    userID,
		SessionID: sessionID,
		Role:      role,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   AuthTokenIssuer,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(AccessTokenLifetime),
			),
		},
	}

	if err := claims.Validate(); err != nil {
		return nil, fmt.Errorf("validate access claims: %w", err)
	}

	return claims, nil
}

func (i *Issuer) IssueRefresh(userID uuid.UUID, sessionID string, lifetime domain.SessionLifetime) (*RefreshClaims, error) {
	now := time.Now()

	claims := &RefreshClaims{
		UserID:    userID,
		SessionID: sessionID,

		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:   AuthTokenIssuer,
			IssuedAt: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(
				now.Add(time.Duration(lifetime)),
			),
		},
	}

	if err := claims.Validate(); err != nil {
		return nil, fmt.Errorf("validate refresh claims: %w", err)
	}

	return claims, nil
}

func (i *Issuer) SignAccess(claims *AccessClaims) (string, error) {
	return sign(i.secret, i.signingMethod, claims)
}

func (i *Issuer) SignRefresh(claims *RefreshClaims) (string, error) {
	return sign(i.secret, i.signingMethod, claims)
}

func (i *Issuer) ParseAccess(tokenString string) (*AccessClaims, error) {
	return parse[AccessClaims](i.secret, tokenString)
}

func (i *Issuer) ParseRefresh(tokenString string) (*RefreshClaims, error) {
	return parse[RefreshClaims](i.secret, tokenString)
}
