package requests_mongo_repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/repository/mongo/pool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) GetRequest(ctx context.Context, id uuid.UUID) (domain.Request, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if id == uuid.Nil {
		return domain.Request{}, fmt.Errorf("`ID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	filter := bson.D{{Key: "_id", Value: id.String()}}

	var model RequestModel
	if err := r.pool.Requests().FindOne(ctx, filter).Decode(&model); err != nil {
		err = core_mongo_pool.MapErrors(err)
		if errors.Is(err, core_mongo_pool.ErrNotFound) {
			return domain.Request{}, fmt.Errorf("request with id='%s': %w", id, core_errors.ErrNotFound)
		}
		return domain.Request{}, fmt.Errorf("mongo find one: %w", err)
	}

	return modelToDomain(model), nil
}
