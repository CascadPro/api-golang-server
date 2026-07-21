package test_postgres_repository

import core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/repository/postgres/pool"

type TestRepository struct {
	pool core_postgres_pool.Pool
}

func NewTestRepository(pool core_postgres_pool.Pool) *TestRepository {
	return &TestRepository{
		pool: pool,
	}
}
