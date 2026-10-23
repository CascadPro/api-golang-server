package core_postgres_outbox

import (
	"context"
	"fmt"
)

func (r *Repository) Cleanup(ctx context.Context) (int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	result, err := r.pool.Exec(ctx, cleanupQuery)
	if err != nil {
		return 0, fmt.Errorf("cleanup processed outbox events: %w", err)
	}

	return result.RowsAffected(), nil
}
