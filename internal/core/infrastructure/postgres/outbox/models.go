package core_postgres_outbox

import (
	"encoding/json"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type Event struct {
	ID uuid.UUID

	Type        domain.EventType
	AggregateID *uuid.UUID
	Payload     json.RawMessage

	Attempts  int
	LastError *string

	ProcessedAt *time.Time
	CreatedAt   time.Time
}

func domainEventFromModel(model Event) domain.OutboxEvent {
	return domain.OutboxEvent{
		ID:          model.ID,
		Type:        model.Type,
		AggregateID: model.AggregateID,
		Payload:     model.Payload,
		Attempts:    model.Attempts,
		LastError:   model.LastError,
		ProcessedAt: model.ProcessedAt,
		CreatedAt:   model.CreatedAt,
	}
}

func domainEventsFromModels(models []Event) []domain.OutboxEvent {
	events := make([]domain.OutboxEvent, len(models))
	for i, model := range models {
		events[i] = domainEventFromModel(model)
	}

	return events
}
