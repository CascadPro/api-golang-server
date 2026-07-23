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
	if err := validateCommonClaims(&t.RegisteredClaims, t.SessionID, t.UserID); err != nil {
		return fmt.Errorf("access token validation: %w", err)
	}

	roles := make(map[domain.UserRole]struct{}, len(domain.Roles))
	for _, r := range domain.Roles {
		roles[r] = struct{}{}
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
	if err := validateCommonClaims(&t.RegisteredClaims, t.SessionID, t.UserID); err != nil {
		return fmt.Errorf("refresh token validation: %w", err)
	}

	return nil
}

func validateCommonClaims(rc *jwt.RegisteredClaims, sessionID string, userID uuid.UUID) error {
	issuedAt, err := rc.GetIssuedAt()
	if err != nil {
		return fmt.Errorf("`IssuedAt` is empty: %w", core_errors.ErrInvalidArgument)
	}

	issuer, err := rc.GetIssuer()
	if err != nil {
		return fmt.Errorf("`Issuer` is empty: %w", core_errors.ErrInvalidArgument)
	}
	if issuer != AuthTokenIssuer {
		return fmt.Errorf("`Issuer` is invalid: %w", core_errors.ErrInvalidArgument)
	}

	if rc.ExpiresAt == nil {
		return fmt.Errorf("`ExpiresAt` is missing: %w", core_errors.ErrInvalidArgument)
	}
	if rc.ExpiresAt.Before(issuedAt.Time) {
		return fmt.Errorf("token is expired: %w", core_errors.ErrUnauthorized)
	}

	if sessionID == "" {
		return fmt.Errorf("`SessionID` is empty: %w", core_errors.ErrInvalidArgument)
	}

	if _, err := uuid.Parse(fmt.Sprint(userID)); err != nil {
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
