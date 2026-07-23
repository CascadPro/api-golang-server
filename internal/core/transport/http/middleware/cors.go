package core_http_middleware

import (
	"net/http"
	"strings"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			config, err := core_config.NewConfig()
			if err != nil {
				config.AllowedOrigins = ""
			}

			allowedOrigins := map[string]struct{}{}
			allowedMethods := []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"}

			for allowedOrigin := range strings.SplitSeq(config.AllowedOrigins, ",") {
				allowedOrigins[strings.TrimSpace(allowedOrigin)] = struct{}{}
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
