package core_postgres_outbox

import (
	"context"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

func (r *Repository) GetPendingEvents(
	ctx context.Context,
	workerID uuid.UUID,
	limit int,
	ttl time.Duration,
) ([]domain.OutboxEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	if workerID == uuid.Nil {
		return nil, fmt.Errorf("`workerID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	rows, err := r.pool.Query(ctx, getPendingQuery, limit, workerID, int64(ttl.Seconds()))
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var events []Event

	for rows.Next() {
		var event Event

		if err := rows.Scan(
			&event.ID,
			&event.Type,
			&event.AggregateID,
			&event.Payload,
			&event.Attempts,
			&event.LastError,
			&event.LockedAt,
			&event.LockedBy,
			&event.ProcessedAt,
			&event.CreatedAt,
		); err != nil {
			return nil, fmt.Errorf("scan error: %w", err)
		}

		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iterator error: %w", err)
	}

	return domainEventsFromModels(events), nil
}
