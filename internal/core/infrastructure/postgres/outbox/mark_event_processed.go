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

	query := `
		UPDATE infrastructure.outbox_events
		SET
			processed_at = $3,
			locked_at = NULL,
			locked_by = NULL,
			last_error = NULL
		WHERE (id = $1 AND locked_by = $2 AND processed_at IS NULL);
	`

	now := time.Now()

	cmd, err := r.pool.Exec(ctx, query, id, workerID, now)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("outbox event %s is not owned by worker %s", id, workerID)
	}

	return nil
}
