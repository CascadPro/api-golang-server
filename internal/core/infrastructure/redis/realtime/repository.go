package core_redis_realtime

import (
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/pool"
)

const ChannelKey = "cascade:realtime"

type Publisher struct {
	pool core_redis_pool.Pool
}

func NewPublisher(pool core_redis_pool.Pool) *Publisher {
	return &Publisher{
		pool: pool,
	}
}
