package core_http_middleware

import (
	"context"
	"net/http"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_i18n "github.com/CascadePro/api-golang-server/internal/core/i18n"
)

func Locale() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			locale := core_i18n.Resolve(r.Header.Get("Accept-Language"))

			ctx := context.WithValue(r.Context(), core_context.CtxKeyLocale, locale)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
