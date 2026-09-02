package requests_mongo_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/mongo/pool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

func (r *Repository) PatchRequestStatus(
	ctx context.Context,
	id uuid.UUID,
	userID uuid.UUID,
	status domain.RequestStatus,
) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if id == uuid.Nil {
		return fmt.Errorf("`ID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if userID == uuid.Nil {
		return fmt.Errorf("`UserID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if status == domain.RequestStatusNil {
		return fmt.Errorf("`Status` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	request, err := r.GetRequest(ctx, id)
	if err != nil {
		return fmt.Errorf("get request: %w", err)
	}

	request.Status = status

	err = request.Validate()
	if err != nil {
		request.Status = domain.RequestStatusWaiting
	}

	now, versionInc := time.Now(), request.Version+1

	update := make([]bson.E, 0, 4)
	update = append(update, bson.E{Key: "v", Value: versionInc}, bson.E{Key: "updated_at", Value: now})

	update = append(update, bson.E{Key: "status", Value: request.Status})
	update = append(update, bson.E{Key: "status_by", Value: userID.String()})

	filter := bson.D{{Key: "_id", Value: id.String()}, {Key: "v", Value: request.Version}}

	if err := defaultPatch(ctx, r.pool.Requests(), filter, update); err != nil {
		return err
	}

	return err
}

func defaultPatch(ctx context.Context, coll *mongo.Collection, filter bson.D, update []bson.E) error {
	opts := options.UpdateOne().SetUpsert(false)

	_, err := coll.UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}}, opts)
	if err != nil {
		err = core_mongo_pool.MapErrors(err)
		if errors.Is(err, core_mongo_pool.ErrNotFound) {
			return fmt.Errorf("request has concurrently accessed: %w", core_errors.ErrConflict)
		}
		return fmt.Errorf("mongo update one: %w", err)
	}

	return nil
}
