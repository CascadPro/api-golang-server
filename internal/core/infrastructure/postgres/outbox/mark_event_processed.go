package core_postgres_outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) MarkEventProcessed(ctx context.Context, id uuid.UUID, workerID uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	now := time.Now()

	cmd, err := r.pool.Exec(ctx, markProcessedQuery, id, workerID, now)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("outbox event %s is not owned by worker %s", id, workerID)
	}

	return nil
}
