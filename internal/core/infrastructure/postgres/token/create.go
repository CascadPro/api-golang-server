package core_postgres_token

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	"github.com/google/uuid"
)

// Create token, fields to create: "Token", "Type", "ExpiresAt", "UserID"
func (r *Repository) Create(ctx context.Context, token domain.Token) (domain.Token, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	id, err := uuid.NewV7()
	if err != nil {
		return domain.Token{}, fmt.Errorf("generate ID: %w", err)
	}

	row := r.pool.QueryRow(ctx, createQuery, id, token.Token, token.Type, token.ExpiresAt, token.UserID)

	var model TokenModel
	if err := row.Scan(&model.ID, &model.Version, &model.Token, &model.Type, &model.ExpiresAt); err != nil {
		err = core_postgres_pool.MapErrors(err)

		if errors.Is(err, core_postgres_pool.ErrViolatesForeignKey) {
			return domain.Token{}, fmt.Errorf("%v: user with id=%s: %w", err, token.UserID, core_errors.ErrNotFound)
		}

		return domain.Token{}, fmt.Errorf("scan error: %w", err)
	}

	return domainFromModel(model), nil
}
