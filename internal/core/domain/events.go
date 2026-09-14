package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type EventType string

const (
	EventTypeMediaDeleteFile    = EventType("media.delete_file")
	EventTypeMediaAvatarProcess = EventType("media.process_avatar")
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

type EventTypeMediaDeleteFilePayload struct {
	Tag    FileTag `json:"tag"`
	FileID string  `json:"file_id"`
}

func NewEventTypeMediaDeleteFilePayload(tag FileTag, fileID string) EventTypeMediaDeleteFilePayload {
	return EventTypeMediaDeleteFilePayload{
		Tag:    tag,
		FileID: fileID,
	}
}

type EventTypeMediaAvatarProcessPayload struct {
	FileID string `json:"file_id"`
}

func NewEventTypeMediaAvatarProcessPayload(fileID string) EventTypeMediaAvatarProcessPayload {
	return EventTypeMediaAvatarProcessPayload{
		FileID: fileID,
	}
}
