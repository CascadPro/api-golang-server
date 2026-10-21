package core_http_middleware

import (
	"net/http"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_http_utils "github.com/CascadePro/api-golang-server/internal/core/transport/http/utils"
	"github.com/google/uuid"
	"go.uber.org/zap"
)

const requestIDHeader = "X-Request-ID"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			if requestID == "" {
				requestID = uuid.NewString()
			}
			r.Header.Set(requestIDHeader, requestID)
			w.Header().Set(requestIDHeader, requestID)
			ctx := core_context.WithRequestID(r.Context(), requestID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIDHeader)
			l := log.With(zap.String("request_id", requestID), zap.String("url", r.URL.String()))
			ctx := core_logger.ToContext(r.Context(), l)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			locale := core_context.Locale(ctx)
			responseHandler := core_http_response.NewResponseHandler(log, locale, w)
			defer func() {
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle http request got unexpected panic")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			rw := core_http_response.NewResponseWriter(w)
			before := time.Now()
			log.Debug(">>> incoming http request", zap.String("method", r.Method), zap.Time("time", time.Now().UTC()))
			next.ServeHTTP(rw, r)
			log.Debug("<<< done http request", zap.Int("status_code", rw.GetStatusCode()), zap.Duration("latency", time.Since(before)))
		})
	}
}

// IP resolves the remote client address using only explicitly trusted proxy ranges.
// In dev/offline environments callers can still use the configured fixed IP behavior
// by choosing it at the composition root; this middleware does not trust headers by default.
func IP(trustedProxyCIDRs ...string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)

			ip, err := core_http_utils.ClientIPFromRequest(r, trustedProxyCIDRs)
			if err != nil {
				log.Warn("failed to determine client IP", zap.Error(err))
				http.Error(w, "unable to determine client IP", http.StatusBadRequest)
				return
			}

			log.Debug("resolved client IP", zap.String("ip", ip.String()))

			ctx = core_context.WithClientIP(ctx, ip)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
