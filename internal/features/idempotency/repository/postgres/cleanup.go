package idempotency_postgres_repository

import (
	"context"
	"fmt"
)

func (r *Repository) CleanupExpired(ctx context.Context) (int64, error) {
	result, err := r.pool.Exec(ctx, cleanupExpiredQuery)
	if err != nil {
		return 0, fmt.Errorf("cleanup expired idempotency records: %w", err)
	}

	return result.RowsAffected(), nil
}
