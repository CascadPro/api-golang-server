package core_ws_hub

import (
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_ws_client "github.com/CascadePro/api-golang-server/internal/core/transport/ws/client"
	"github.com/google/uuid"
)

func (h *Hub) Broadcast(event domain.RealtimeEvent) {
	h.mu.RLock()

	clients := make([]*core_ws_client.Client, 0, len(h.clients))

	for client := range h.clients {
		clients = append(clients, client)
	}

	h.mu.RUnlock()

	for _, client := range clients {
		if client.SendEvent(event) {
			continue
		}

		h.Unregister(client)
	}
}

func (h *Hub) SendToUser(userID uuid.UUID, event domain.RealtimeEvent) {
	h.mu.RLock()

	clients := make([]*core_ws_client.Client, 0, len(h.byUser[userID]))

	for client := range h.byUser[userID] {
		clients = append(clients, client)
	}

	h.mu.RUnlock()

	for _, client := range clients {
		var ok bool

		switch event.Type {
		case domain.RealtimeEventSessionRevoked:
			ok = client.RevokeSession(event)

		case domain.RealtimeEventSessionsRevoked:
			ok = client.RevokeAllSessions(event)

		default:
			ok = client.SendEvent(event)
		}

		if !ok {
			h.Unregister(client)
		}
	}
}

func (h *Hub) SendToSession(sessionID string, event domain.RealtimeEvent) {
	h.mu.RLock()

	clients := make([]*core_ws_client.Client, 0, len(h.bySession[sessionID]))

	for client := range h.bySession[sessionID] {
		clients = append(clients, client)
	}

	h.mu.RUnlock()

	for _, client := range clients {
		var ok bool

		ok = client.SendEvent(event)

		if !ok {
			h.Unregister(client)
		}
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
