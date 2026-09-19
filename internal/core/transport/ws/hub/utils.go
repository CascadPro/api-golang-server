package core_ws_hub

import (
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

func newPresenceOnline(userID uuid.UUID, sessionID string) domain.RealtimeEvent {
	event := domain.NewRealtimeEvent(
		domain.RealtimeEventPresenceOnline,
		nil,
		nil,
		domain.NewRealtimeEventPresenceData(userID, sessionID),
	)

	return event
}

func newPresenceOffline(userID uuid.UUID, sessionID string) domain.RealtimeEvent {
	event := domain.NewRealtimeEvent(
		domain.RealtimeEventPresenceOffline,
		nil,
		nil,
		domain.NewRealtimeEventPresenceData(userID, sessionID),
	)

	return event
}
