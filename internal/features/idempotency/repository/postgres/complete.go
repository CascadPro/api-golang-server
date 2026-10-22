package idempotency_postgres_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
)

func (r *Repository) Complete(ctx context.Context, record domain.IdempotencyRecord, status int, body []byte, headers http.Header) error {
	return completeWith(ctx, r.pool, record, status, body, headers)
}

func (r *Repository) CompleteTx(ctx context.Context, tx core_postgres_pool.Tx, record domain.IdempotencyRecord, status int, body []byte, headers http.Header) error {
	return completeWith(ctx, tx, record, status, body, headers)
}

func completeWith(ctx context.Context, db core_postgres_pool.Querier, record domain.IdempotencyRecord, status int, body []byte, headers http.Header) error {
	if body == nil {
		body = []byte{}
	}
	if headers == nil {
		headers = make(http.Header)
	}
	encodedHeaders, err := json.Marshal(headers)
	if err != nil {
		return fmt.Errorf("marshal response headers: %w", err)
	}

	result, err := db.Exec(ctx, completeQuery, record.Scope, record.Operation, record.KeyHash, status, body, encodedHeaders)
	if err != nil {
		return fmt.Errorf("complete idempotency record: %w", err)
	}
	if result.RowsAffected() != 1 {
		return fmt.Errorf("idempotency record was not in processing state")
	}
	return nil
}
