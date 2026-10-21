package idempotency_postgres_repository

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/jackc/pgx/v5"
)

func scanRecord(row pgx.Row, scope, operation string, keyHash []byte) (domain.IdempotencyRecord, error) {
	var record domain.IdempotencyRecord
	var status *int16
	var body, encodedHeaders []byte

	err := row.Scan(
		&record.RequestHash,
		&record.Status,
		&status,
		&body,
		&encodedHeaders,
	)
	if err != nil {
		return domain.IdempotencyRecord{}, err
	}

	record.Scope = scope
	record.Operation = operation
	record.KeyHash = keyHash
	record.ResponseBody = body
	record.ResponseHeaders = make(http.Header)

	if status != nil {
		record.ResponseStatus = int(*status)
	}

	if len(encodedHeaders) > 0 {
		if err := json.Unmarshal(encodedHeaders, &record.ResponseHeaders); err != nil {
			return domain.IdempotencyRecord{}, fmt.Errorf("decode response headers: %w", err)
		}
	}

	return record, nil
}
