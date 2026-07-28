package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
)

type Router struct {
	*http.ServeMux
	prefix     string
	middleware []core_http_middleware.Middleware
}

func NewRouter(path string, middleware ...core_http_middleware.Middleware) *Router {
	return &Router{
		ServeMux:   http.NewServeMux(),
		prefix:     path,
		middleware: middleware,
	}
}

func (r *Router) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *Router) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(r, r.middleware...)
}
