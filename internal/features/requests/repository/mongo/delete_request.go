package requests_mongo_repository

import (
	"context"
	"errors"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/repository/mongo/pool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (r *Repository) DeleteRequest(ctx context.Context, id uuid.UUID) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if id == uuid.Nil {
		return fmt.Errorf("`ID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	filter := bson.D{{Key: "_id", Value: id.String()}}

	if _, err := r.pool.Requests().DeleteOne(ctx, filter); err != nil {
		err = core_mongo_pool.MapErrors(err)
		if errors.Is(err, core_mongo_pool.ErrNotFound) {
			return fmt.Errorf("request with id='%s': %w", id, core_errors.ErrNotFound)
		}
		return fmt.Errorf("mongo delete one: %w", err)
	}

	return nil
}
