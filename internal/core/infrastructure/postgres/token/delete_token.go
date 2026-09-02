package core_postgres_token

import (
	"context"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (r *Repository) DeleteToken(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if id == uuid.Nil {
		return fmt.Errorf("`id` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	query := `
		DELETE FROM base.tokens
		WHERE id = $1;
	`

	cmd, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("token with id=%s: %w", id, core_errors.ErrNotFound)
	}

	return nil
}
