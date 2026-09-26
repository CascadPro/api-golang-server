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
	mu      sync.RWMutex
	closing bool

	clients   map[*core_ws_client.Client]struct{}
	byUser    map[uuid.UUID]map[*core_ws_client.Client]struct{}
	bySession map[string]map[*core_ws_client.Client]struct{}

	ctx context.Context

	logger *core_logger.Logger

	publisher core_redis_realtime.PublisherMethods
	presence  presence_service.ServiceMethods
}

func NewHub(
	ctx context.Context,
	publisher core_redis_realtime.PublisherMethods,
	presence presence_service.ServiceMethods,
	logger *core_logger.Logger,
) *Hub {
	return &Hub{
		clients:   make(map[*core_ws_client.Client]struct{}),
		byUser:    make(map[uuid.UUID]map[*core_ws_client.Client]struct{}),
		bySession: make(map[string]map[*core_ws_client.Client]struct{}),

		ctx: ctx,

		logger: logger,

		publisher: publisher,
		presence:  presence,
	}
}

func (h *Hub) Register(client *core_ws_client.Client) error {
	h.mu.Lock()

	if h.closing {
		h.mu.Unlock()
		client.Close()

		return fmt.Errorf("websocket hub is shutting down")
	}

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

	result, err := h.presence.Register(h.ctx, client.UserID, client.SessionID, client.ID)
	if err != nil {
		h.removeClient(client)

		return fmt.Errorf("register presence: %w", err)
	}

	if !result.WasSessionOnline {
		event := newPresenceOnline(client.UserID, client.SessionID)

		if err := h.publisher.Publish(h.ctx, event); err != nil {
			h.logger.Error("publish presence online", zap.Error(err))
		}
	}

	return nil
}

func (h *Hub) removeClient(client *core_ws_client.Client) bool {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, exists := h.clients[client]; !exists {
		return false
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

	return true
}

func (h *Hub) Unregister(client *core_ws_client.Client) {
	if !h.removeClient(client) {
		return
	}

	client.Close()

	result, err := h.presence.Unregister(h.ctx, client.UserID, client.SessionID, client.ID)
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

		if err := h.publisher.Publish(h.ctx, event); err != nil {
			h.logger.Error("publish presence offline", zap.Error(err))
		}
	}
}

func (h *Hub) Shutdown(ctx context.Context) error {
	h.mu.Lock()

	if h.closing {
		h.mu.Unlock()
		return nil
	}

	h.closing = true

	clients := make([]*core_ws_client.Client, 0, len(h.clients))
	for client := range h.clients {
		clients = append(clients, client)
	}

	h.mu.Unlock()

	for _, client := range clients {
		client.Shutdown()
	}

	var firstErr error

	for _, client := range clients {
		if err := client.Wait(ctx); err != nil && firstErr == nil {
			firstErr = err
		}
	}

	return firstErr
}
