package core_ws_client

import (
	"context"
	"sync"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_ws_conn "github.com/CascadePro/api-golang-server/internal/core/transport/ws/conn"
	core_ws_dispatcher "github.com/CascadePro/api-golang-server/internal/core/transport/ws/dispatcher"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	pingPeriod     = 10 * time.Second
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
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

	dispatcher *core_ws_dispatcher.Dispatcher

	closeOnce sync.Once

	pumps sync.WaitGroup
}

func NewClient(
	conn *core_ws_conn.Conn,
	userID uuid.UUID,
	sessionID string,
	dispatcher *core_ws_dispatcher.Dispatcher,
) *Client {
	return &Client{
		ID:        uuid.NewString(),
		UserID:    userID,
		SessionID: sessionID,

		conn: conn,
		send: make(chan ChannelMessage, sendBufferSize),
		done: make(chan struct{}),

		dispatcher: dispatcher,
	}
}

func (c *Client) Serve(unregister func(*Client)) {
	c.pumps.Add(2)

	go func() {
		defer c.pumps.Done()
		c.readPump(unregister)
	}()

	go func() {
		defer c.pumps.Done()
		c.writePump(unregister)
	}()
}

func (c *Client) Shutdown() {
	c.closeOnce.Do(func() {
		closeBody := websocket.FormatCloseMessage(websocket.CloseGoingAway, "server shutting down")

		_ = c.conn.WriteControl(
			websocket.CloseMessage,
			closeBody,
			time.Now().Add(writeWait),
		)

		close(c.done)

		_ = c.conn.Close()

		c.conn.CancelContext()
	})
}

func (c *Client) Wait(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		defer close(done)
		c.pumps.Wait()
	}()

	select {
	case <-done:
		return nil

	case <-ctx.Done():
		return ctx.Err()
	}
}

func (c *Client) Close() {
	c.closeOnce.Do(func() {
		close(c.done)

		_ = c.conn.Close()

		c.conn.CancelContext()
	})
}

func (c *Client) logMarshalError(id string, t domain.RealtimeEventType, err error) {
	log := core_logger.FromContext(c.conn.Context())

	log.Error(
		"marshal websocket event",
		zap.String("event_id", id),
		zap.String("event_type", string(t)),
		zap.Error(err),
	)
}

func (c *Client) logError(msg string, err error) {
	log := core_logger.FromContext(c.conn.Context())

	log.Error(
		msg,
		zap.Error(err),
		zap.String("client_id", c.ID),
		zap.String("session_id", c.SessionID),
		zap.String("user_id", c.UserID.String()),
	)
}
