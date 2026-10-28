package idempotency_service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"net/http"
	"time"

	core_context "github.com/CascadePro/api-golang-server/internal/core/context"
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_http_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/http/middleware"
	core_http_response "github.com/CascadePro/api-golang-server/internal/core/transport/http/response"
	core_http_utils "github.com/CascadePro/api-golang-server/internal/core/transport/http/utils"
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
				err = fmt.Errorf("%w: %v", idempotency_error.ErrInvalidIdempotencyKey, err)
				responseHandler.ErrorResponse(err, "failed to get idempotency key")
				return
			}

			fingerprint, err := requestFingerprint(r)
			if err != nil {
				responseHandler.ErrorResponse(err, "failed to fingerprint request")
				return
			}

			keyHash := sha256.Sum256([]byte(key.String()))
			operation := r.Method + ":" + r.URL.EscapedPath()
			scope := userID.String()

			reserveCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
			record, acquired, err := s.repository.Reserve(reserveCtx, scope, operation, keyHash[:], fingerprint, time.Now().Add(s.ttl))

			cancel()

			if err != nil {
				responseHandler.ErrorResponse(err, "failed to reserve idempotency key")
				return
			}

			if !acquired {
				if !equalBytes(record.RequestHash, fingerprint) {
					responseHandler.ErrorResponse(idempotency_error.ErrIdempotencyKeyReused, "idempotency key reused with a different request")
					return
				}

				switch record.Status {
				case domain.IdempotencyRecordStatusProcessing:
					rw.Header().Set("Retry-After", "1")
					responseHandler.ErrorResponse(idempotency_error.ErrIdempotencyInProgress, "idempotency operation is still processing")
					return

				case domain.IdempotencyRecordStatusCompleted:
					replayResponse(rw, record)
					return

				default:
					err = fmt.Errorf("unknown idempotency status %q", record.Status)
					responseHandler.ErrorResponse(err, "invalid idempotency record status")
					return
				}
			}

			completeFn := func(completeCtx context.Context, tx core_postgres_pool.Tx, status int, body []byte, headers http.Header) error {
				return s.repository.CompleteTx(completeCtx, tx, record, status, body, headers)
			}

			ctx = withTransactionState(ctx, record, completeFn)
			r = r.WithContext(ctx)

			recorder := core_http_utils.NewRecorder()

			next.ServeHTTP(recorder, r)

			recStatus, recBody, recHeaders := recorder.GetStatus(), recorder.GetBody(), recorder.Header()

			if !transactionCommitted(ctx) {
				completeCtx, completeCancel := context.WithTimeout(context.Background(), 5*time.Second)
				defer completeCancel()

				if err := s.repository.Complete(completeCtx, record, recStatus, recBody.Bytes(), safeHeaders(recHeaders)); err != nil {
					responseHandler.ErrorResponse(err, "failed to save idempotency response")
					return
				}
			}

			copyHeaders(rw.Header(), recHeaders)

			rw.WriteHeader(recStatus)

			if recStatus != http.StatusNoContent && recStatus != http.StatusNotModified {
				_, _ = rw.Write(recBody.Bytes())
			}
		})
	}
}
