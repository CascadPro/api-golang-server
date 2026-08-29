package core_http_middleware

import (
	"context"
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_http_utils "github.com/CascadePro/api-golang-server/internal/core/transport/http/utils"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func Authorization(issuer core_jwt_security.AccessTokenVerifier, rolesArr ...domain.UserRole) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewResponseHandler(log, rw)

			accessToken, err := core_http_utils.ParseAuthorizationHeader(r, issuer)
			if err != nil {
				responseHandler.ErrorResponse(err, "parse authorization header")
				return
			}

			if len(rolesArr) != 0 {
				if err := core_validation.ValidateArray(rolesArr, accessToken.Role); err != nil {
					responseHandler.ErrorResponse(core_errors.ErrForbidden, "you don't have enough rights")
					return
				}
			}

			ctx = context.WithValue(ctx, core_context.CtxKeyUserID, accessToken.UserID.String())
			ctx = context.WithValue(ctx, core_context.CtxKeyUserRole, accessToken.Role)
			ctx = context.WithValue(ctx, core_context.CtxKeySessionID, accessToken.SessionID)

			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
