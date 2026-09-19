package core_ws_hub

import (
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

func (h *Hub) Broadcast(event domain.RealtimeEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		client.SendEvent(event)
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, event domain.RealtimeEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.byUser[userID] {
		client.SendEvent(event)
	}
}

func (h *Hub) SendToSession(sessionID string, event domain.RealtimeEvent) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.bySession[sessionID] {
		if event.Type == domain.RealtimeEventSessionRevoked {
			client.Kick(event)
			continue
		}

		client.SendEvent(event)
	}
}

func (h *Hub) HandleEvent(event domain.RealtimeEvent) {
	if len(event.Recipients) > 0 {
		for _, userID := range event.Recipients {
			h.SendToUser(userID, event)
		}

		return
	}

	if event.SessionID != nil {
		h.SendToSession(*event.SessionID, event)
		return
	}

	if event.UserID != nil {
		h.SendToUser(*event.UserID, event)
		return
	}

	h.Broadcast(event)
}
