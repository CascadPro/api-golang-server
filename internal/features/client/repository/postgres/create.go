package client_postgres_repository

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	"github.com/google/uuid"
)

// Create client, fields to create. "Company", "Contacts"
func (r *Repository) Create(ctx context.Context, client domain.Client) (domain.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	id, err := uuid.NewV7()
	if err != nil {
		return domain.Client{}, fmt.Errorf("generate ID: %w", err)
	}

	row := r.pool.QueryRow(ctx, createQuery, id, client.Company, client.Contacts)

	var model ClientModel
	if err := row.Scan(&model.ID, &model.Version, &model.Company, &model.Contacts, &model.CreatedAt); err != nil {
		return domain.Client{}, fmt.Errorf("scan error: %w", core_postgres_pool.MapErrors(err))
	}

	return domainFromModel(model), nil
}
