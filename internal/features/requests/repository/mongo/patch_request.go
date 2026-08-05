package requests_mongo_repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/repository/mongo/pool"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// PatchRequest updates selected fields of an existing request using $set.
// Patched fields: "Title", "Status", "Origin", "Docs", "WorkTypes", "Geography", "Deadline", "Contract".
func (r *Repository) PatchRequest(ctx context.Context, id uuid.UUID, request domain.Request) error {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if id == uuid.Nil {
		return fmt.Errorf("`ID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if request.Version < 1 {
		return fmt.Errorf("`Version` can't be uninitialized: %w", core_errors.ErrInvalidArgument)
	}

	model := domainToModel(request)

	now := time.Now()
	versionInc := model.Version + 1

	update := make([]bson.E, 0, 12)
	update = append(update, bson.E{Key: "v", Value: versionInc}, bson.E{Key: "updated_at", Value: now})

	if model.Title != "" {
		update = append(update, bson.E{Key: "title", Value: model.Title})
	}
	if model.Status != domain.RequestStatusNil {
		update = append(update, bson.E{Key: "status", Value: model.Status})
	}
	if len(model.Origin) != 0 {
		update = append(update, bson.E{Key: "origin", Value: model.Origin})
	}
	if len(model.WorkTypes) != 0 {
		update = append(update, bson.E{Key: "work_types", Value: model.WorkTypes})
	}
	if len(model.Geography) != 0 {
		update = append(update, bson.E{Key: "geo_desc", Value: model.Geography})
	}
	if model.Deadline != nil {
		update = append(update, bson.E{Key: "deadline", Value: model.Deadline})
	}
	if model.ContractDocID != nil {
		update = append(update, bson.E{Key: "contract", Value: model.ContractDocID})
	}
	if model.StatusBy != nil {
		update = append(update, bson.E{Key: "status_by", Value: model.StatusBy})
	}
	if model.ClientID != nil {
		update = append(update, bson.E{Key: "client", Value: model.ClientID})
	}
	if model.Docs.TechTaskDocID != nil || model.Docs.ProjectDocID != nil || model.Docs.SpecificationDocID != nil {
		update = append(update, bson.E{Key: "docs", Value: model.Docs})
	}

	idStr := id.String()
	filter := bson.D{{Key: "_id", Value: idStr}, {Key: "v", Value: request.Version}}

	_, err := r.pool.Requests().UpdateOne(ctx, filter, bson.D{{Key: "$set", Value: update}})
	if err != nil {
		err = core_mongo_pool.MapErrors(err)
		if errors.Is(err, core_mongo_pool.ErrNotFound) {
			return fmt.Errorf("request with id='%s' has concurrently accessed: %w", id, core_errors.ErrConflict)
		}
		return fmt.Errorf("mongo update one: %w", err)
	}

	return nil
}
