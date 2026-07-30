package core_http_middleware

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func Authorization(rolesArr ...domain.UserRole) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewResponseHandler(log, rw)

			bearer := r.Header.Get("Authorization")
			if bearer == "" {
				err := fmt.Errorf("authorization header empty: %w", core_errors.ErrUnauthorized)
				responseHandler.ErrorResponse(err, "Header `Authorization` is empty")
				return
			}

			parts := strings.SplitN(bearer, " ", 2)
			if len(parts) != 2 || strings.TrimSpace(parts[0]) != "Bearer" {
				err := fmt.Errorf("invalid `Authorization` header format")
				responseHandler.ErrorResponse(core_errors.ErrUnauthorized, err.Error())
				return
			}
			tokenString := strings.TrimSpace(parts[1])

			claims, err := core_utils.GetJwtTokenClaims[core_utils.JwtAccessClaims](tokenString)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to read token")
				return
			}

			if len(rolesArr) != 0 {
				if err := core_validation.ValidateArray(rolesArr, claims.Role); err != nil {
					responseHandler.ErrorResponse(core_errors.ErrForbidden, "you don't have enough rights")
					return
				}
			}

			ctx = context.WithValue(ctx, core_context.CtxKeyUserID, claims.UserID.String())
			ctx = context.WithValue(ctx, core_context.CtxKeyUserRole, claims.Role)
			ctx = context.WithValue(ctx, core_context.CtxKeySessionID, claims.SessionID)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
