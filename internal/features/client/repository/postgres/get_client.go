package client_postgres_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	"github.com/google/uuid"
)

func (r *Repository) GetClient(ctx context.Context, id uuid.UUID) (domain.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if id == uuid.Nil {
		return domain.Client{}, fmt.Errorf("`ID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	query := `
		SELECT id, version, company, contacts, created_at
		FROM base.clients
		WHERE (id = $1)
		LIMIT 1;
	`

	row := r.pool.QueryRow(ctx, query, id)

	var model ClientModel
	if err := row.Scan(&model.ID, &model.Version, &model.Company, &model.Contacts, &model.CreatedAt); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrNoRows) {
			return domain.Client{}, fmt.Errorf("%v: client with id='%s': %w", err, id, core_errors.ErrNotFound)
		}

		return domain.Client{}, fmt.Errorf("scan error: %w", err)
	}

	return domainFromModel(model), nil
}
