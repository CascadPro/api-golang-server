package requests_mongo_repository

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/mongo/pool"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func (r *Repository) GetStats(ctx context.Context) (domain.RequestStats, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	stats := domain.RequestStats{
		ByStatus: make(map[domain.RequestStatus]int64, len(domain.RequestStatuses)),
	}

	for _, status := range domain.RequestStatuses {
		stats.ByStatus[status] = 0
	}

	pipeline := mongo.Pipeline{
		bson.D{{Key: "$facet", Value: bson.D{
			{
				Key: "total",
				Value: mongo.Pipeline{
					bson.D{{Key: "$count", Value: "count"}},
				},
			},
			{
				Key: "statuses",
				Value: mongo.Pipeline{
					bson.D{{Key: "$group", Value: bson.D{
						{Key: "_id", Value: "$status"},
						{Key: "count", Value: bson.D{
							{Key: "$sum", Value: 1},
						}},
					}}},
				},
			},
		}}},
	}

	cursor, err := r.pool.Requests().Aggregate(ctx, pipeline)
	if err != nil {
		return stats, fmt.Errorf(
			"mongo aggregate request stats: %w",
			core_mongo_pool.MapErrors(err),
		)
	}
	defer cursor.Close(ctx)

	var results []RequestStatsAggregation
	if err := cursor.All(ctx, &results); err != nil {
		return stats, fmt.Errorf(
			"decode request stats: %w",
			core_mongo_pool.MapErrors(err),
		)
	}

	if len(results) == 0 {
		return stats, nil
	}

	result := results[0]

	if len(result.Total) > 0 {
		stats.Total = result.Total[0].Count
	}

	for _, item := range result.Statuses {
		if _, exists := stats.ByStatus[item.Status]; exists {
			stats.ByStatus[item.Status] = item.Count
		}
	}

	return stats, nil
}
