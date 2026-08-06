package core_http_middleware

import (
	"net/http"
	"strings"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
	core_http "github.com/CascadePro/api-golang-server/internal/core/transport/http"
)

func CORS() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			cfg, err := core_config.NewConfig()
			if err != nil {
				cfg.AllowedOrigins = ""
			}

			allowedOrigins := map[string]struct{}{}
			allowedMethods := make([]string, len(core_http.AllowedMethods))

			for allowedOrigin := range strings.SplitSeq(cfg.AllowedOrigins, ",") {
				allowedOrigins[strings.TrimSpace(allowedOrigin)] = struct{}{}
			}

			for i, method := range core_http.AllowedMethods {
				allowedMethods[i] = string(method)
			}

			origin := r.Header.Get("Origin")

			if _, ok := allowedOrigins[origin]; ok {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Set("Access-Control-Allow-Methods", strings.Join(allowedMethods, ", "))
				w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			}

			if r.Method == string(core_http.MethodOptions) {
				w.WriteHeader(http.StatusOK)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
