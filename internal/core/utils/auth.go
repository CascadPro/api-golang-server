package core_utils

import (
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

const (
	RefreshTokenLifetime time.Duration = time.Hour * 24 * 30
	AccessTokenLifetime  time.Duration = time.Minute * 15
	AuthTokenIssuer      string        = "cascade-pro-api-server"
)

type JwtAccessClaims struct {
	UserID    uuid.UUID       `json:"uid"`
	SessionID string          `json:"sid"`
	Role      domain.UserRole `json:"role"`

	jwt.RegisteredClaims
}

func (t *JwtAccessClaims) Validate() error {
	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return fmt.Errorf("`IssuedAt` is empty: %w", core_errors.ErrInvalidArgument)
	}

	issuer, err := t.GetIssuer()
	if err != nil {
		return fmt.Errorf("`Issuer` is empty: %w", core_errors.ErrInvalidArgument)
	}

	if issuer != AuthTokenIssuer {
		return fmt.Errorf("`Issuer` is invalid: %w", core_errors.ErrInvalidArgument)
	}

	if t.ExpiresAt.Before(issuedAt.Time) {
		return fmt.Errorf("token is expired: %w", core_errors.ErrUnauthorized)
	}

	if t.SessionID == "" {
		return fmt.Errorf("`SessionID` is empty: %w", core_errors.ErrInvalidArgument)
	}

	if _, err := uuid.Parse(fmt.Sprint(t.UserID)); err != nil {
		return fmt.Errorf("`UserID` is invalid: %w", core_errors.ErrInvalidArgument)
	}

	roles := map[domain.UserRole]struct{}{}
	for _, role := range domain.Roles {
		roles[role] = struct{}{}
	}

	if _, ok := roles[t.Role]; !ok {
		return fmt.Errorf("`Role` is invalid: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type JwtRefreshClaims struct {
	UserID    uuid.UUID `json:"uid"`
	SessionID string    `json:"sid"`

	jwt.RegisteredClaims
}

func (t *JwtRefreshClaims) Validate() error {
	issuedAt, err := t.GetIssuedAt()
	if err != nil {
		return fmt.Errorf("`IssuedAt` is empty: %w", core_errors.ErrInvalidArgument)
	}

	issuer, err := t.GetIssuer()
	if err != nil {
		return fmt.Errorf("`Issuer` is empty: %w", core_errors.ErrInvalidArgument)
	}

	if issuer != AuthTokenIssuer {
		return fmt.Errorf("`Issuer` is invalid: %w", core_errors.ErrInvalidArgument)
	}

	if t.ExpiresAt.Before(issuedAt.Time) {
		return fmt.Errorf("token is expired: %w", core_errors.ErrUnauthorized)
	}

	if t.SessionID == "" {
		return fmt.Errorf("`SessionID` is empty: %w", core_errors.ErrInvalidArgument)
	}

	if _, err := uuid.Parse(fmt.Sprint(t.UserID)); err != nil {
		return fmt.Errorf("`UserID` is invalid: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func IssueTokens(userID uuid.UUID, sessionID string, role domain.UserRole) (*JwtAccessClaims, *JwtRefreshClaims, error) {
	accessToken := JwtAccessClaims{
		UserID:    userID,
		SessionID: sessionID,
		Role:      role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    AuthTokenIssuer,
		},
	}

	if err := accessToken.Validate(); err != nil {
		return nil, nil, fmt.Errorf("issue tokens: %w", err)
	}

	refreshToken := JwtRefreshClaims{
		UserID:    userID,
		SessionID: sessionID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshTokenLifetime)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    AuthTokenIssuer,
		},
	}

	if err := refreshToken.Validate(); err != nil {
		return nil, nil, fmt.Errorf("issue tokens: %w", err)
	}

	return &accessToken, &refreshToken, nil
}
