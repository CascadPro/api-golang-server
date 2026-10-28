package requests_mongo_repository

import (
	core_mongo_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/mongo/pool"
)

type Repository struct {
	pool core_mongo_pool.Pool
}

func NewRepository(pool core_mongo_pool.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
