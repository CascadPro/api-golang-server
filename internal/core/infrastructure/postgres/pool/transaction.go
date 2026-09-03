package core_postgres_pool

import (
	"context"
)

type Tx interface {
	Querier

	Commit(ctx context.Context) error

	Rollback(ctx context.Context) error
}
