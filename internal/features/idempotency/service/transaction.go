package idempotency_service

import (
	"context"
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
)

type transactionStateKey struct{}

type transactionState struct {
	record    domain.IdempotencyRecord
	complete  func(context.Context, core_postgres_pool.Tx, int, []byte, http.Header) error
	committed bool
}

// CompleteInTransaction stores the response inside the business transaction.
// Call it before committing the transaction.
func CompleteInTransaction(
	ctx context.Context,
	tx core_postgres_pool.Tx,
	status int,
	body []byte,
	headers http.Header,
) error {
	state, ok := ctx.Value(transactionStateKey{}).(*transactionState)
	if !ok || state == nil || state.complete == nil {
		return nil
	}

	return state.complete(ctx, tx, status, body, headers)
}

// MarkTransactionCommitted must only be called after tx.Commit succeeds.
func MarkTransactionCommitted(ctx context.Context) {
	state, ok := ctx.Value(transactionStateKey{}).(*transactionState)
	if ok && state != nil {
		state.committed = true
	}
}

func withTransactionState(
	ctx context.Context,
	record domain.IdempotencyRecord,
	complete func(context.Context, core_postgres_pool.Tx, int, []byte, http.Header) error,
) context.Context {
	state := &transactionState{
		record:   record,
		complete: complete,
	}

	return context.WithValue(ctx, transactionStateKey{}, state)
}

func transactionCommitted(ctx context.Context) bool {
	state, ok := ctx.Value(transactionStateKey{}).(*transactionState)

	return ok && state != nil && state.committed
}
