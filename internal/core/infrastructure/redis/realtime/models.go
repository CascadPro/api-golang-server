package core_redis_realtime

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type EventModel struct {
	ID         string                   `json:"id"`
	Type       domain.RealtimeEventType `json:"type"`
	UserID     *uuid.UUID               `json:"user_id,omitempty"`
	SessionID  *string                  `json:"session_id,omitempty"`
	Recipients []uuid.UUID              `json:"recipients,omitempty"`
	Data       any                      `json:"data,omitempty"`
	CreatedAt  time.Time                `json:"timestamp"`
}

func modelToDomain(event EventModel) domain.RealtimeEvent {
	return domain.RealtimeEvent{
		ID:         event.ID,
		Type:       event.Type,
		UserID:     event.UserID,
		SessionID:  event.SessionID,
		Recipients: event.Recipients,
		Data:       event.Data,
		CreatedAt:  event.CreatedAt,
	}
}

func domainToModel(event domain.RealtimeEvent) EventModel {
	return EventModel{
		ID:         event.ID,
		Type:       event.Type,
		UserID:     event.UserID,
		SessionID:  event.SessionID,
		Recipients: event.Recipients,
		Data:       event.Data,
		CreatedAt:  event.CreatedAt,
	}
}
