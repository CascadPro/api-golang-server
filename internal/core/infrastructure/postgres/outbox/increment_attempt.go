package core_postgres_outbox

import (
	"context"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (r *Repository) IncrementAttempt(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		UPDATE infrastructure.outbox_events
		SET attempts = attempts + 1
		WHERE (id = $1 AND processed_at IS NULL);
	`

	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("outbox event with id=%s: %w", id, core_errors.ErrNotFound)
	}

	return nil
}
