package core_http_server

import (
	"fmt"
	"net/http"

	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
)

type ApiVersion string

const (
	ApiVersion1 = ApiVersion("v1")
)

type ApiVersionRouter struct {
	*http.ServeMux
	version    ApiVersion
	middleware []core_http_middleware.Middleware
}

func NewApiVersionRouter(version ApiVersion, middleware ...core_http_middleware.Middleware) *ApiVersionRouter {
	return &ApiVersionRouter{
		ServeMux:   http.NewServeMux(),
		version:    version,
		middleware: middleware,
	}
}

func (r *ApiVersionRouter) RegisterRoutes(routes ...Route) {
	for _, route := range routes {
		pattern := fmt.Sprintf("%s %s", route.Method, route.Path)

		r.Handle(pattern, route.WithMiddleware())
	}
}

func (r *ApiVersionRouter) RegisterRouters(routers ...*Router) {
	for _, router := range routers {
		r.Handle(router.prefix+"/", http.StripPrefix(router.prefix, router.WithMiddleware()))
	}
}

func (r *ApiVersionRouter) WithMiddleware() http.Handler {
	return core_http_middleware.ChainMiddleware(r, r.middleware...)
}
