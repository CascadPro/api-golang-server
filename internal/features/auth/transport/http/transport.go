package auth_transport_http

import (
	"context"
	"net"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_request "github.com/CascadePro/api-golang-server/internal/core/transport/http/request"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
)

type HttpHandler struct {
	authService AuthService
	tokenIssuer core_jwt_security.IssuerMethods
	rateLimiter core_http_middleware.HTTPRateLimitMiddleware
}

type AuthService interface {
	CreateRegisterToken(context.Context, domain.User) (domain.Token, error)

	Register(context.Context, domain.UserPatch, string) error

	Login(context.Context, domain.User, net.IP, *core_http_request.UserAgent) (string, string, error)

	GetNewTokens(context.Context, string) (string, error)

	Logout(context.Context, string) error
}

func NewHttpHandler(
	authService AuthService,
	tokenIssuer core_jwt_security.IssuerMethods,
	rateLimiter core_http_middleware.HTTPRateLimitMiddleware,
) *HttpHandler {
	return &HttpHandler{
		authService: authService,
		tokenIssuer: tokenIssuer,
		rateLimiter: rateLimiter,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	roles := []domain.UserRole{domain.RoleAdmin, domain.RoleDirector}

	return []core_http_server.Route{
		{
			Method:     core_http.MethodPost,
			Path:       "/register/token",
			Handler:    h.CreateRegisterToken,
			Middleware: []core_http_middleware.Middleware{core_http_middleware.Authorization(h.tokenIssuer, roles...)},
		},
		{
			Method:     core_http.MethodPost,
			Path:       "/register",
			Handler:    h.Register,
			Middleware: []core_http_middleware.Middleware{h.rateLimiter.Middleware()},
		},
		{
			Method:     core_http.MethodPost,
			Path:       "/login",
			Handler:    h.Login,
			Middleware: []core_http_middleware.Middleware{h.rateLimiter.Middleware()},
		},
		{
			Method:     core_http.MethodPost,
			Path:       "/logout",
			Handler:    h.Logout,
			Middleware: []core_http_middleware.Middleware{h.rateLimiter.Middleware()},
		},
		{
			Method:  core_http.MethodGet,
			Path:    "/login/refresh",
			Handler: h.GetNewTokens,
		},
	}
}
