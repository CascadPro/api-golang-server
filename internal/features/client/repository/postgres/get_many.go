package client_postgres_repository

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

const (
	defaultLimit = 20
	defaultPage  = 1
)

func (r *Repository) GetMany(
	ctx context.Context,
	page, limit int64,
	searchQuery *string,
) ([]domain.Client, int64, error) {
	if limit <= 0 {
		limit = defaultLimit
	}

	if page <= 0 {
		page = defaultPage
	}

	query := ""
	if searchQuery != nil {
		query = *searchQuery
	}

	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	offset := (page - 1) * limit

	rows, err := r.pool.Query(ctx, getManyQuery, query, limit, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query: %w", err)
	}

	var models []ClientModel
	var totalCount int64

	for rows.Next() {
		var model ClientModel

		if err := rows.Scan(
			&model.ID,
			&model.Version,
			&model.Company,
			&model.Contacts,
			&model.CreatedAt,
			&totalCount,
		); err != nil {
			rows.Close()
			return nil, 0, fmt.Errorf("scan: %w", err)
		}

		models = append(models, model)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return nil, 0, fmt.Errorf("rows: %w", err)
	}
	rows.Close()

	if len(models) == 0 {
		err := r.pool.QueryRow(ctx, getCountQuery, query).Scan(&totalCount)
		if err != nil {
			return nil, 0, fmt.Errorf("count: %w", err)
		}
	}

	return domainsFromModels(models), totalCount, nil
}
