package core_mongo_pool

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Pool interface {
	OpTimeout() time.Duration

	Close(ctx context.Context)
}

type ConnectionPool struct {
	client  *mongo.Client
	pool    *mongo.Database
	timeout time.Duration

	Collections Collections
}

type Collections struct {
	Requests *mongo.Collection
	Chapters *mongo.Collection
}

func New(ctx context.Context, cfg Config) (*ConnectionPool, error) {
	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)

	opts := options.Client().
		ApplyURI(uri).
		SetMaxPoolSize(100).
		SetMinPoolSize(5).
		SetMaxConnIdleTime(5 * time.Minute)

	client, err := mongo.Connect(opts)
	if err != nil {
		return nil, fmt.Errorf("connect to MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		_ = client.Disconnect(ctx)
		return nil, fmt.Errorf("ping MongoDB: %w", err)
	}

	pool := client.Database(cfg.Database)

	return &ConnectionPool{
		client:  client,
		pool:    pool,
		timeout: cfg.Timeout,
		Collections: Collections{
			Requests: pool.Collection("requests"),
			Chapters: pool.Collection("chapters"),
		},
	}, nil
}

func (p *ConnectionPool) OpTimeout() time.Duration {
	return p.timeout
}

func (p *ConnectionPool) Close(ctx context.Context) {
	_ = p.client.Disconnect(ctx)
}
