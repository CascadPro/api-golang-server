package core_postgres_outbox

import (
	"context"
	"fmt"

	"github.com/google/uuid"
)

func (r *Repository) MarkEventFailed(ctx context.Context, id uuid.UUID, workerID uuid.UUID, lastError error) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	cmd, err := r.pool.Exec(ctx, markFailedQuery, id, workerID, lastError.Error())
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("outbox event %s is not owned by worker %s", id, workerID)
	}

	return nil
}
