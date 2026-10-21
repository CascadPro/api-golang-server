package idempotency_service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	idempotency_error "github.com/CascadePro/api-golang-server/internal/features/idempotency/errors"
	"github.com/google/uuid"
)

func (s *Service) Middleware() core_http_middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			locale := core_context.Locale(ctx)

			responseHandler := core_http_response.NewResponseHandler(log, locale, rw)

			userID, err := core_context.UserID(ctx)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to get userID from context")
				return
			}

			key, err := uuid.Parse(r.Header.Get(HeaderKey))
			if err != nil {
				err = fmt.Errorf("%v: %w", err, idempotency_error.ErrInvalidIdempotencyKey)
				responseHandler.ErrorResponse(err, "failed to get idempotency key")
				return
			}

			fingerprint, err := requestFingerprint(r)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to decode request body")
				return
			}

			keyHash := sha256.Sum256([]byte(key.String()))
			operation := r.Method + ":" + r.URL.EscapedPath()
			scope := userID.String()

			ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
			defer cancel()

			record, acquired, err := s.repository.Reserve(
				ctx,
				scope,
				operation,
				keyHash[:],
				fingerprint,
				time.Now().Add(s.ttl),
			)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to check idempotency key")
				return
			}

			if !acquired {
				if !equalBytes(record.RequestHash, fingerprint) {
					err = fmt.Errorf("%v: %w", err, idempotency_error.ErrIdempotencyKeyReused)
					responseHandler.ErrorResponse(err, "failed to use idempotency key")
					return
				}

				if record.Status == domain.IdempotencyRecordStatusProcessing {
					rw.Header().Set("Retry-After", "1")
					responseHandler.ErrorResponse(idempotency_error.ErrIdempotencyInProgress, "idempotency now in progress")
					return
				}

				replayResponse(rw, record)
				return
			}

			recorder := newRecorder()
			next.ServeHTTP(recorder, r)

			completeCtx, completeCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer completeCancel()

			if err := s.repository.Complete(
				completeCtx,
				record,
				recorder.status,
				recorder.body.Bytes(),
				safeHeaders(recorder.header),
			); err != nil {
				responseHandler.ErrorResponse(err, "failed to save request body")
				return
			}

			copyHeaders(rw.Header(), recorder.header)

			rw.WriteHeader(recorder.status)
			_, _ = rw.Write(recorder.body.Bytes())
		})
	}
}
