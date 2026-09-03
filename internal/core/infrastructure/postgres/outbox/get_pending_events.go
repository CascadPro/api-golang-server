package core_postgres_outbox

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

func (r *Repository) GetPendingEvents(ctx context.Context, limit int) ([]domain.OutboxEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		SELECT id, aggregate_id, aggregate_type, event_type, payload, attempts, last_error, processed_at, created_at
		FROM infrastructure.outbox_events
		WHERE processed_at IS NULL
		ORDER BY created_at, id
		LIMIT $1
		FOR UPDATE SKIP LOCKED;
	`

	rows, err := r.pool.Query(ctx, query, limit)
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
			&event.CreatedAt,
			&event.ProcessedAt,
			&event.Attempts,
			&event.LastError,
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
