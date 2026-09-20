package core_ws_client

import (
	"sync"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_ws_conn "github.com/CascadePro/api-golang-server/internal/core/transport/ws/conn"
	core_ws_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/ws/middleware"
	presence_redis_repository "github.com/CascadePro/api-golang-server/internal/features/presence/repository/redis"
	presence_service "github.com/CascadePro/api-golang-server/internal/features/presence/service"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = 25 * time.Second
	maxMessageSize = 64 * 1024
	sendBufferSize = 32
)

type Client struct {
	ID        string
	UserID    uuid.UUID
	SessionID string

	conn *core_ws_conn.Conn
	send chan ChannelMessage
	done chan struct{}

	presence presence_service.ServiceMethods

	closeOnce sync.Once
}

func NewClient(
	conn *core_ws_conn.Conn,
	userID uuid.UUID,
	sessionID string,
	presence presence_service.ServiceMethods,
) *Client {
	return &Client{
		ID:        uuid.NewString(),
		UserID:    userID,
		SessionID: sessionID,

		conn: conn,
		send: make(chan ChannelMessage, sendBufferSize),
		done: make(chan struct{}),

		presence: presence,
	}
}

func (c *Client) Serve(unregister func(*Client)) {
	go c.readPump(unregister)
	go c.writePump(unregister)
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.done)

		_ = c.conn.Close()

		c.conn.CancelContext()
	})
}

func (c *Client) readPump(unregister func(*Client)) {
	defer unregister(c)

	c.conn.SetReadLimit(maxMessageSize)

	if err := c.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
		return
	}

	c.conn.SetPongHandler(func(string) error {
		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})

	handler := core_ws_middleware.ChainMiddleware(
		incomingHandler,
		core_ws_middleware.Panic(),
	)

	for {
		if err := handler(c.conn); err != nil {
			return
		}
	}
}

func (c *Client) writePump(unregister func(*Client)) {
	ticker := time.NewTicker(pingPeriod)
	defer ticker.Stop()

	presenceTicker := time.NewTicker(presence_redis_repository.HeartbeatInterval)
	defer presenceTicker.Stop()

	for {
		select {
		case <-c.done:
			return

		case message := <-c.send:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				unregister(c)

				c.conn.CancelContext()

				return
			}

			if err := c.conn.WriteMessage(websocket.TextMessage, message.data); err != nil {
				unregister(c)

				c.conn.CancelContext()

				return
			}

			if message.closeAfter {
				closeBody := websocket.FormatCloseMessage(int(message.closeCode), string(message.closeReason))

				_ = c.conn.WriteControl(websocket.CloseMessage, closeBody, time.Now().Add(writeWait))

				unregister(c)

				c.conn.CancelContext()
				return
			}

		case <-ticker.C:
			if err := c.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
				unregister(c)

				c.conn.CancelContext()

				return
			}

			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				unregister(c)

				c.conn.CancelContext()

				return
			}

		case <-presenceTicker.C:
			if err := c.presence.Heartbeat(
				c.conn.Context(),
				c.UserID,
				c.SessionID,
				c.ID,
			); err != nil {
				unregister(c)

				c.conn.CancelContext()

				return
			}
		}
	}
}

func (c *Client) logError(id string, t domain.RealtimeEventType, err error) {
	log := core_logger.FromContext(c.conn.Context())

	log.Error(
		"marshal websocket event",
		zap.String("event_id", id),
		zap.String("event_type", string(t)),
		zap.Error(err),
	)
}
