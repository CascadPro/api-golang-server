package requests_mongo_repository

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/mongo/pool"
	"github.com/google/uuid"
)

func (r *Repository) CreateRequest(ctx context.Context, request domain.Request) (uuid.UUID, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if err := request.Validate(); err != nil {
		return uuid.Nil, fmt.Errorf("validate request: %w", err)
	}

	id, err := uuid.NewV7()
	if err != nil {
		return uuid.Nil, fmt.Errorf("generate ID: %w", err)
	}

	model := domainToModel(request)

	model.ID = id.String()
	model.Version = 1

	model.CalcRequiredEmptyFields()

	if _, err := r.pool.Requests().InsertOne(ctx, model); err != nil {
		return uuid.Nil, fmt.Errorf("mongo insert one: %w", core_mongo_pool.MapErrors(err))
	}

	return uuid.MustParse(model.ID), nil
}
