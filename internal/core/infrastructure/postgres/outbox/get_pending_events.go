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

	query := `
		WITH events AS (
			SELECT id
			FROM infrastructure.outbox_events
			WHERE processed_at IS NULL AND (
				locked_at IS NULL
				OR locked_at < $3 - ($4 * INTERVAL '1 second')
			)
			ORDER BY created_at, id
			FOR UPDATE SKIP LOCKED
			LIMIT $1
		)
		UPDATE infrastructure.outbox_events AS e
		SET locked_at = $3, locked_by = $2, attempts = attempts + 1
		FROM events
		WHERE e.id = events.id
		RETURNING
			e.id,
			e.type,
			e.aggregate_id,
			e.payload,
			e.attempts,
			e.last_error,
			e.locked_at,
			e.locked_by,
			e.processed_at,
			e.created_at;
	`

	now := time.Now()

	rows, err := r.pool.Query(ctx, query, limit, workerID, now, int64(ttl.Seconds()))
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
