package core_jwt_security

import (
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type AccessClaims struct {
	UserID    uuid.UUID       `json:"uid"`
	SessionID string          `json:"sid"`
	Role      domain.UserRole `json:"role"`

	jwt.RegisteredClaims
}

func (c *AccessClaims) Validate() error {
	if err := validateCommonClaims(&c.RegisteredClaims, c.SessionID, c.UserID); err != nil {
		return fmt.Errorf("access token validation: %w", err)
	}

	roles := make(
		map[domain.UserRole]struct{},
		len(domain.Roles),
	)

	for _, role := range domain.Roles {
		roles[role] = struct{}{}
	}

	if _, ok := roles[c.Role]; !ok {
		return fmt.Errorf("`Role` is invalid: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type RefreshClaims struct {
	UserID    uuid.UUID `json:"uid"`
	SessionID string    `json:"sid"`

	jwt.RegisteredClaims
}

func (c *RefreshClaims) Validate() error {
	if err := validateCommonClaims(&c.RegisteredClaims, c.SessionID, c.UserID); err != nil {
		return fmt.Errorf("refresh token validation: %w", err)
	}

	return nil
}

var (
	_ core_validation.Validatable = (*AccessClaims)(nil)
	_ core_validation.Validatable = (*RefreshClaims)(nil)
)
