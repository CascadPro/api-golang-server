package idempotency_postgres_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func (r *Repository) Reserve(ctx context.Context, scope, operation string, keyHash, requestHash []byte, expiresAt time.Time) (domain.IdempotencyRecord, bool, error) {
	record, err := scanRecord(
		r.pool.QueryRow(ctx, reserveQuery, scope, operation, keyHash, requestHash, expiresAt),
		scope,
		operation,
		keyHash,
	)
	if err == nil {
		return record, true, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return domain.IdempotencyRecord{}, false, fmt.Errorf("reserve idempotency key: %w", err)
	}

	record, err = scanRecord(
		r.pool.QueryRow(ctx, scanRecordQuery, scope, operation, keyHash),
		scope,
		operation,
		keyHash,
	)
	if err != nil {
		return domain.IdempotencyRecord{}, false, fmt.Errorf("get idempotency record: %w", err)
	}

	return record, false, nil
}
