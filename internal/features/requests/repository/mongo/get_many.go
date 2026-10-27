package requests_mongo_repository

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/mongo/pool"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

const (
	defaultLimit = 10
	defaultPage  = 1
)

func (r *Repository) GetRequests(
	ctx context.Context,
	page, limit int,
	filter domain.RequestQueryFilter,
) ([]domain.Request, int64, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if limit <= 0 {
		limit = defaultLimit
	}
	if page <= 0 {
		page = defaultPage
	}

	if err := filter.Validate(); err != nil {
		return nil, 0, fmt.Errorf("validate `filter`: %w", err)
	}

	mongoFilter := bson.M{}
	if len(filter.Status) > 0 {
		mongoFilter["status"] = bson.M{"$in": filter.Status}
	}

	total, err := r.pool.Requests().CountDocuments(ctx, mongoFilter)
	if err != nil {
		return nil, 0, fmt.Errorf("mongo count documents: %w", core_mongo_pool.MapErrors(err))
	}

	if total <= 0 {
		return []domain.Request{}, 0, nil
	}

	sortOrder := -1
	if filter.Sort == domain.SortTypeOldest {
		sortOrder = 1
	}

	opts := options.Find().
		SetSkip(int64((page - 1) * limit)).SetLimit(int64(limit)).
		SetSort(bson.D{{Key: "created_at", Value: sortOrder}})

	cursor, err := r.pool.Requests().Find(ctx, mongoFilter, opts)
	if err != nil {
		return nil, 0, fmt.Errorf("mongo find: %w", core_mongo_pool.MapErrors(err))
	}
	defer cursor.Close(ctx)

	var models []RequestModel

	for cursor.Next(ctx) {
		var model RequestModel

		if err := cursor.Decode(&model); err != nil {
			return nil, 0, fmt.Errorf("decode request model: %w", core_mongo_pool.MapErrors(err))
		}

		models = append(models, model)
	}

	if err := cursor.Err(); err != nil {
		return nil, 0, fmt.Errorf("iterate request cursor: %w", core_mongo_pool.MapErrors(err))
	}

	requests, err := modelsToDomains(models)
	if err != nil {
		return nil, 0, fmt.Errorf("convert models to domains: %w", err)
	}

	return requests, total, nil
}
