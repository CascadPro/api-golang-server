package media_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

// PatchFile, fields to update: "Filename", "Size", "Deleted", "DeletedAt"
func (r *Repository) PatchFile(ctx context.Context, fileID string, file domain.File) (domain.File, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if err := core_validation.ValidateID(fileID, domain.FileIDByteLength); err != nil {
		return domain.File{}, fmt.Errorf("validate file id: %w", err)
	}

	query := `
		UPDATE media.files
		SET (version = version + 1, filename = $1, size = $2, deleted = $3, deleted_at = $4)
		WHERE (id = $5 AND version = $6)
		RETURNING id, version, tag, filename, mime_type, size, deleted, deleted_at, created_at;
	`

	row := r.pool.QueryRow(ctx, query, file.Filename, file.Size, file.Deleted, file.DeletedAt, fileID, file.Version)

	var model FileModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Tag,
		&model.Filename,
		&model.MimeType,
		&model.Size,
		&model.Deleted,
		&model.DeletedAt,
		&model.CreatedAt,
	); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.File{}, fmt.Errorf(
				"file with id=%s concurrently accessed: %w",
				fileID, core_errors.ErrConflict,
			)
		}

		return domain.File{}, fmt.Errorf("scan error: %w", err)
	}

	return domainFromModel(model), nil
}
