package media_postgres_repository

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (r *Repository) DeleteFile(ctx context.Context, fileID string) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return fmt.Errorf("validate file id: %w", err)
	}

	query := `
		DELETE FROM media.files
		WHERE id = $1;
	`

	cmd, err := r.pool.Exec(ctx, query, fileID)
	if err != nil {
		return fmt.Errorf("exec query: %w", err)
	}
	if cmd.RowsAffected() == 0 {
		return fmt.Errorf("file with id=%s: %w", fileID, core_errors.ErrNotFound)
	}

	return nil
}
