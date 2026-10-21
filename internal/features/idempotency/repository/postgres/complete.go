package idempotency_postgres_repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func (r *Repository) Complete(ctx context.Context, record domain.IdempotencyRecord, status int, body []byte, headers http.Header) error {
	encodedHeaders, err := json.Marshal(headers)
	if err != nil {
		return fmt.Errorf("marshal response headers: %w", err)
	}

	result, err := r.pool.Exec(
		ctx,
		completeQuery,
		record.Scope,
		record.Operation,
		record.KeyHash,
		status,
		body,
		encodedHeaders,
	)
	if err != nil {
		return fmt.Errorf("complete idempotency record: %w", err)
	}

	if result.RowsAffected() != 1 {
		return fmt.Errorf("idempotency record was not in processing state")
	}

	return nil
}
