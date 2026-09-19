package core_ws_hub

import (
	"context"
	"fmt"
	"sync"

	"github.com/google/uuid"
	"go.uber.org/zap"

	core_redis_realtime "github.com/CascadePro/api-golang-server/internal/core/infrastructure/redis/realtime"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_ws_client "github.com/CascadePro/api-golang-server/internal/core/transport/ws/client"
	presence_service "github.com/CascadePro/api-golang-server/internal/features/presence/service"
)

type Hub struct {
	mu sync.RWMutex

	clients   map[*core_ws_client.Client]struct{}
	byUser    map[uuid.UUID]map[*core_ws_client.Client]struct{}
	bySession map[string]map[*core_ws_client.Client]struct{}

	logger *core_logger.Logger

	publisher core_redis_realtime.PublisherMethods
	presence  presence_service.ServiceMethods
}

func NewHub(publisher core_redis_realtime.PublisherMethods, presence presence_service.ServiceMethods, logger *core_logger.Logger) *Hub {
	return &Hub{
		clients:   make(map[*core_ws_client.Client]struct{}),
		byUser:    make(map[uuid.UUID]map[*core_ws_client.Client]struct{}),
		bySession: make(map[string]map[*core_ws_client.Client]struct{}),

		logger: logger,

		publisher: publisher,
		presence:  presence,
	}
}

func (h *Hub) Register(client *core_ws_client.Client) error {
	h.mu.Lock()

	h.clients[client] = struct{}{}

	if h.byUser[client.UserID] == nil {
		h.byUser[client.UserID] =
			make(map[*core_ws_client.Client]struct{})
	}

	h.byUser[client.UserID][client] = struct{}{}

	if h.bySession[client.SessionID] == nil {
		h.bySession[client.SessionID] =
			make(map[*core_ws_client.Client]struct{})
	}

	h.bySession[client.SessionID][client] = struct{}{}

	h.mu.Unlock()

	becameOnline, err := h.presence.Register(context.Background(), client.UserID, client.SessionID, client.ID)
	if err != nil {
		h.Unregister(client)

		return fmt.Errorf("register presence: %w", err)
	}

	if becameOnline {
		event := newPresenceOnline(client.UserID, client.SessionID)
		_ = h.publisher.Publish(context.Background(), event)
	}

	return nil
}

func (h *Hub) Unregister(client *core_ws_client.Client) {
	h.mu.Lock()

	if _, exists := h.clients[client]; !exists {
		h.mu.Unlock()
		return
	}

	delete(h.clients, client)

	userClients := h.byUser[client.UserID]

	delete(userClients, client)

	if len(userClients) == 0 {
		delete(h.byUser, client.UserID)
	}

	sessionClients := h.bySession[client.SessionID]

	delete(sessionClients, client)

	if len(sessionClients) == 0 {
		delete(h.bySession, client.SessionID)
	}

	h.mu.Unlock()

	client.Close()

	result, err := h.presence.Unregister(context.Background(), client.UserID, client.SessionID, client.ID)
	if err != nil {
		h.logger.Error(
			"failed to unregister websocket presence",
			zap.Error(err),
			zap.String("client_id", client.ID),
			zap.String("session_id", client.SessionID),
			zap.String("user_id", client.UserID.String()),
		)

		return
	}

	if result.SessionOffline {
		event := newPresenceOffline(client.UserID, client.SessionID)

		_ = h.publisher.Publish(context.Background(), event)
	}

	if result.UserOffline {
		// event := newUserPresenceOffline(client.UserID)
		// _ = h.publisher.Publish(context.Background(), event)
	}
}
