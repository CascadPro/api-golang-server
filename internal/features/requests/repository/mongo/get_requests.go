package requests_mongo_repository

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/mongo/pool"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"go.uber.org/zap"
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
	log := core_logger.FromContext(ctx)
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
			log.Error("mongo cursor decode model", zap.Error(core_mongo_pool.MapErrors(err)))
			continue
		}

		models = append(models, model)
	}
	if err := cursor.Err(); err != nil {
		return nil, total, fmt.Errorf("mongo cursor error: %w", core_mongo_pool.MapErrors(err))
	}

	return modelsToDomains(models), total, nil
}
