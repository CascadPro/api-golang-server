package core_http_utils

import (
	"fmt"
	"net/http"
	"strings"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	jwt "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
)

func ParseAuthorizationHeader(r *http.Request, issuer jwt.AccessTokenVerifier) (*jwt.AccessClaims, error) {
	bearer := r.Header.Get("Authorization")
	if bearer == "" {
		return nil, fmt.Errorf("`Authorization` header can't be empty: %w", core_errors.ErrUnauthorized)
	}

	parts := strings.SplitN(bearer, " ", 2)
	if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Bearer" {
		return nil, fmt.Errorf("`Authorization` header has invalid format: %w", core_errors.ErrInvalidArgument)
	}

	tokenString := strings.TrimSpace(parts[1])
	if tokenString == "" {
		return nil, fmt.Errorf("access token string can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	claims, err := issuer.ParseAccess(tokenString)
	if err != nil {
		return nil, fmt.Errorf("parse access token claims: %w", err)
	}

	return claims, nil
}
