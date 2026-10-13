package core_http_middleware

import (
	"net/http"
	"strings"

	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
)

func CORS(allowedOrigins string) Middleware {
	origins := make(map[string]struct{})

	for _, origin := range strings.Split(allowedOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins[origin] = struct{}{}
		}
	}

	methods := make([]string, len(core_http.AllowedMethods))
	for i, method := range core_http.AllowedMethods {
		methods[i] = string(method)
	}

	allowedMethods := strings.Join(methods, ", ")

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			origin := r.Header.Get("Origin")

			if _, ok := origins[origin]; origin != "" && ok {
				w.Header().Set(
					"Access-Control-Allow-Origin",
					origin,
				)
				w.Header().Set(
					"Access-Control-Allow-Methods",
					allowedMethods,
				)
				w.Header().Set(
					"Access-Control-Allow-Headers",
					"Content-Type, Authorization",
				)
				w.Header().Set(
					"Access-Control-Allow-Credentials",
					"true",
				)
				w.Header().Add("Vary", "Origin")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
