package media_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	media_errors "github.com/CascadePro/api-golang-server/internal/features/media/errors"
)

func (r *Repository) PatchPlaceholder(
	ctx context.Context,
	fileID string,
	version int64,
	placeholder []byte,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return fmt.Errorf("validate file id: %w", err)
	}
	if placeholder == nil {
		return fmt.Errorf("placeholder is nil: %w", core_errors.ErrInvalidArgument)
	}

	query := `
	  UPDATE media.files
		SET version = version + 1, placeholder = $1
		WHERE (id = $2 AND version = $3);
	`

	if _, err := r.pool.Exec(ctx, query, placeholder, fileID, version); err != nil {
		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return fmt.Errorf(
				"file with id=%s concurrently accessed: %w",
				fileID, media_errors.ErrFileNotFound,
			)
		}

		return fmt.Errorf("scan error: %w", err)
	}

	return nil
}
