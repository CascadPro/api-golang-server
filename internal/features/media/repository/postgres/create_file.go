package media_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"
	core_utils "github.com/CascadePro/api-golang-server/internal/core/utils"
)

func (r *Repository) CreateFile(ctx context.Context, file domain.File) (domain.File, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO media.files (id, tag, filename, content_type, size)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, version, tag, filename, content_type, size, created_at;
	`

	id, err := core_utils.GenerateID(domain.FileIDByteLength)
	if err != nil {
		return domain.File{}, fmt.Errorf("generate id: %w", err)
	}

	row := r.pool.QueryRow(ctx, query, id, file.Tag, file.Filename, file.ContentType, file.Size)

	var model FileModel
	if err := row.Scan(
		&model.ID,
		&model.Version,
		&model.Tag,
		&model.Filename,
		&model.ContentType,
		&model.Size,
		&model.CreatedAt,
	); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrViolatesUniqueConstraint) {
			return domain.File{}, fmt.Errorf(
				"%v: file with id=%s is already exists: %w",
				err, id, core_errors.ErrInvalidArgument,
			)
		}

		return domain.File{}, fmt.Errorf("scan error: %w", err)
	}

	return domainFromModel(model), nil
}
