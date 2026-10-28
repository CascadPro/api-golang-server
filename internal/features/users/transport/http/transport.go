package users_transport_http

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_server "github.com/CascadePro/api-golang-server/internal/core/transport/http/server"
	"github.com/google/uuid"
)

type HttpHandler struct {
	usersService UsersService
	rateLimiter  core_http_middleware.HTTPRateLimitMiddleware
	idempotency  core_http_middleware.Middleware
}

type UsersService interface {
	GetCurrentUser(ctx context.Context, userID uuid.UUID) (domain.User, []byte, error)

	UpdateAvatar(ctx context.Context, userID uuid.UUID, file *domain.File, content []byte) error

	DeleteAvatar(ctx context.Context, userID uuid.UUID) error
}

func NewHttpHandler(
	usersService UsersService,
	rateLimiter core_http_middleware.HTTPRateLimitMiddleware,
	idempotency core_http_middleware.Middleware,
) *HttpHandler {
	return &HttpHandler{
		usersService: usersService,
		rateLimiter:  rateLimiter,
		idempotency:  idempotency,
	}
}

func (h *HttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  core_http.MethodGet,
			Path:    "/my",
			Handler: h.GetCurrentUser,
		},
		{
			Method:  core_http.MethodPatch,
			Path:    "/avatar",
			Handler: h.UpdateAvatar,
			Middleware: []core_http_middleware.Middleware{
				h.rateLimiter.Middleware(),
				h.idempotency,
				core_http_middleware.Media(),
			},
		},
		{
			Method:  core_http.MethodDelete,
			Path:    "/avatar",
			Handler: h.DeleteAvatar,
			Middleware: []core_http_middleware.Middleware{
				h.rateLimiter.Middleware(),
				h.idempotency,
			},
		},
	}
}
