package core_jwt_security

import (
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

func validateCommonClaims(claims *jwt.RegisteredClaims, sessionID string, userID uuid.UUID) error {
	issuedAt, err := claims.GetIssuedAt()
	if err != nil {
		return fmt.Errorf("`IssuedAt` is empty: %w", core_errors.ErrInvalidArgument)
	}

	issuer, err := claims.GetIssuer()
	if err != nil {
		return fmt.Errorf("`Issuer` is empty: %w", core_errors.ErrInvalidArgument)
	}

	if issuer != AuthTokenIssuer {
		return fmt.Errorf("`Issuer` is invalid: %w", core_errors.ErrInvalidArgument)
	}

	if claims.ExpiresAt == nil {
		return fmt.Errorf("`ExpiresAt` is missing: %w", core_errors.ErrInvalidArgument)
	}

	if claims.ExpiresAt.Before(issuedAt.Time) {
		return fmt.Errorf("token is expired: %w", core_errors.ErrUnauthorized)
	}

	if sessionID == "" {
		return fmt.Errorf("`SessionID` is empty: %w", core_errors.ErrInvalidArgument)
	}

	if userID == uuid.Nil {
		return fmt.Errorf("`UserID` is invalid: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func sign(secret []byte, method jwt.SigningMethod, claims jwt.Claims) (string, error) {
	token := jwt.NewWithClaims(method, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", fmt.Errorf("sign jwt token: %w", err)
	}

	return tokenString, nil
}

func parse[T jwt.Claims](secret []byte, tokenString string) (*T, error) {
	var claims T

	token, err := jwt.ParseWithClaims(
		tokenString,
		any(&claims).(jwt.Claims),
		func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method %v: %w",
					token.Header["alg"],
					core_errors.ErrInvalidArgument,
				)
			}

			return secret, nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse jwt token: %w: %w", err, core_errors.ErrUnauthorized)
	}

	if token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token: %w", core_errors.ErrUnauthorized)
	}

	if validatable, ok := token.Claims.(core_validation.Validatable); ok {
		if err := validatable.Validate(); err != nil {
			return nil, fmt.Errorf("validate jwt claims: %w", err)
		}
	}

	return &claims, nil
}
