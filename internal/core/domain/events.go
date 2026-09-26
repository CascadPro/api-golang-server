package domain

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeMediaDeleteFile    = EventType("media.delete_file")
	EventTypeMediaAvatarProcess = EventType("media.process_avatar")

	EventTypeRealtimePublish = EventType("realtime_events.publish")
)

type OutboxEvent struct {
	ID uuid.UUID

	Type        EventType
	AggregateID *uuid.UUID
	Payload     json.RawMessage

	Attempts  int
	LastError *string

	LockedAt *time.Time
	LockedBy *uuid.UUID

	ProcessedAt *time.Time
	CreatedAt   time.Time
}

func NewOutboxEvent(eventType EventType, aggregateID *uuid.UUID, payload json.RawMessage) OutboxEvent {
	return OutboxEvent{
		ID:          UninitializedUUID,
		Type:        eventType,
		AggregateID: aggregateID,
		Payload:     payload,
		CreatedAt:   time.Now(),
	}
}

type EventMediaDeleteFilePayload struct {
	Tag    FileTag `json:"tag"`
	FileID string  `json:"file_id"`
}

func NewEventMediaDeleteFilePayload(tag FileTag, fileID string) EventMediaDeleteFilePayload {
	return EventMediaDeleteFilePayload{
		Tag:    tag,
		FileID: fileID,
	}
}

type EventMediaAvatarProcessPayload struct {
	FileID string `json:"file_id"`
}

func NewEventMediaAvatarProcessPayload(fileID string) EventMediaAvatarProcessPayload {
	return EventMediaAvatarProcessPayload{
		FileID: fileID,
	}
}

type EventRealtimePublishPayload struct {
	Type      RealtimeEventType `json:"type"`
	UserID    *uuid.UUID        `json:"uid,omitempty"`
	SessionID *string           `json:"sid,omitempty"`
	Data      json.RawMessage   `json:"data,omitempty"`
}

func NewEventRealtimePublishPayload(
	t RealtimeEventType,
	userID *uuid.UUID,
	sessionID *string,
	data any,
) (EventRealtimePublishPayload, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return EventRealtimePublishPayload{}, fmt.Errorf("marshal data: %w", err)
	}

	return EventRealtimePublishPayload{
		Type:      t,
		UserID:    userID,
		SessionID: sessionID,
		Data:      json.RawMessage(raw),
	}, nil
}
