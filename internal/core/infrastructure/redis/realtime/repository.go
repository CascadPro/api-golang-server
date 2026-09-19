package core_redis_realtime

import (
	"context"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
)

const ChannelKey = "cascade:realtime"

type Publisher struct {
	pool core_redis_pool.Pool
}

type PublisherMethods interface {
	Publish(ctx context.Context, event domain.RealtimeEvent) error
	Subscribe(ctx context.Context, handler func(domain.RealtimeEvent)) error
}

func NewPublisher(pool core_redis_pool.Pool) *Publisher {
	return &Publisher{
		pool: pool,
	}
}
