package media_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

func (r *Repository) GetFile(ctx context.Context, fileID string) (domain.File, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return domain.File{}, fmt.Errorf("validate file id: %w", err)
	}

	query := `
		SELECT id, version, tag, filename, content_type, size, deleted, deleted_at, created_at
		FROM media.files
		WHERE (id = $1)
		LIMIT 1;
	`

	row := r.pool.QueryRow(ctx, query, fileID)

	var model FileModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Tag,
		&model.Filename,
		&model.ContentType,
		&model.Size,
		&model.Deleted,
		&model.DeletedAt,
		&model.CreatedAt,
	); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.File{}, fmt.Errorf("%v: file with id=%s: %w", err, fileID, core_errors.ErrNotFound)
		}

		return domain.File{}, fmt.Errorf("scan error: %w", err)
	}

	return domainFromModel(model), nil
}
