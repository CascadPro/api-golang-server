package core_utils

import (
	"fmt"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/golang-jwt/jwt/v5"
)

type validatable interface {
	Validate() error
}

func GenerateJwtToken(method jwt.SigningMethod, claims jwt.Claims) (string, error) {
	config, err := core_config.NewSecretConfig()
	if err != nil {
		return "", fmt.Errorf("get secret config: %w", err)
	}

	token := jwt.NewWithClaims(method, claims)

	tokenString, err := token.SignedString([]byte(config.JwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("sign jwt token: %w", err)
	}

	return tokenString, nil
}

func GetJwtTokenClaims[TP jwt.Claims](tokenString string) (*TP, error) {
	config, err := core_config.NewSecretConfig()
	if err != nil {
		return nil, fmt.Errorf("get secret config: %w", err)
	}

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signature method: %v: %w", t.Header["alg"], core_errors.ErrInvalidArgument)
		}

		return []byte(config.JwtSecretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse jwt token: %w", err)
	}

	if token == nil || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if v, ok := token.Claims.(validatable); ok {
		if err := v.Validate(); err != nil {
			return nil, fmt.Errorf("invalid token: %w", err)
		}
	}

	claims, ok := token.Claims.(TP)
	if !ok {
		return nil, fmt.Errorf("get jwt token payload")
	}

	return &claims, nil
}
