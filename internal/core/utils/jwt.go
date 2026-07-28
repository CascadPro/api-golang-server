package core_utils

import (
	"fmt"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/golang-jwt/jwt/v5"
)

func GenerateJwtToken(method jwt.SigningMethod, claims jwt.Claims) (string, error) {
	cfg, err := core_config.NewSecretConfig()
	if err != nil {
		return "", fmt.Errorf("get secret config: %w", err)
	}

	token := jwt.NewWithClaims(method, claims)

	tokenString, err := token.SignedString([]byte(cfg.JwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("sign jwt token: %w", err)
	}

	return tokenString, nil
}

func GetJwtTokenClaims[TP jwt.Claims](tokenString string) (*TP, error) {
	cfg, err := core_config.NewSecretConfig()
	if err != nil {
		return nil, fmt.Errorf("get secret config: %w", err)
	}

	var claims TP

	token, err := jwt.ParseWithClaims(
		tokenString,
		any(&claims).(jwt.Claims),
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf(
					"unexpected signing method %v: %w",
					t.Header["alg"],
					core_errors.ErrInvalidArgument,
				)
			}
			return []byte(cfg.JwtSecretKey), nil
		},
	)
	if err != nil {
		return nil, fmt.Errorf("parse jwt token: %w", err)
	}

	if token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if v, ok := token.Claims.(core_validation.Validatable); ok {
		if err := v.Validate(); err != nil {
			return nil, fmt.Errorf("invalid token payload: %w", err)
		}
	}

	return &claims, nil
}
