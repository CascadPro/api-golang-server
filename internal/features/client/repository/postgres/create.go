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

	model := modelFromDomain(client)
	model.ID = id

	row := r.pool.QueryRow(ctx, createQuery, model.ID, model.Company, model.Contacts)

	var modelResult ClientModel
	if err := row.Scan(&modelResult.ID, &modelResult.Version, &modelResult.Company, &modelResult.Contacts, &modelResult.CreatedAt); err != nil {
		return domain.Client{}, fmt.Errorf("scan error: %w", core_postgres_pool.MapErrors(err))
	}

	return domainFromModel(modelResult), nil
}