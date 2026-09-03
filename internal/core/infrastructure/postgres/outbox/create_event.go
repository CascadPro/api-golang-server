package core_postgres_outbox

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_postgres_pool "github.com/CascadePro/api-golang-server/internal/core/infrastructure/postgres/pool"
	"github.com/google/uuid"
)

func (r *Repository) CreateEvent(
	ctx context.Context,
	tx core_postgres_pool.Tx,
	event domain.OutboxEvent,
) (domain.OutboxEvent, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := `
		INSERT INTO infrastructure.outbox_events (id, type, aggregate_id, payload)
		VALUES ($1, $2, $3, $4)
		RETURNING id, type, aggregate_id, payload, created_at;
	`

	id, err := uuid.NewV7()
	if err != nil {
		return domain.OutboxEvent{}, fmt.Errorf("generate ID: %w", err)
	}

	if event.Payload == nil {
		event.Payload = json.RawMessage(`{}`)
	}

	row := tx.QueryRow(ctx, query, id, event.Type, event.AggregateID, event.Payload, event.CreatedAt)

	var model Event
	if err := row.Scan(&model.ID, &model.Type, &model.AggregateID, &model.Payload, &model.CreatedAt); err != nil {
		return domain.OutboxEvent{}, fmt.Errorf("scan error: %w", err)
	}

	return domainEventFromModel(model), nil
}
