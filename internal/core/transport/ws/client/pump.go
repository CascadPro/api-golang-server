package core_ws_client

import (
	"time"

	core_ws_middleware "github.com/CascadePro/api-golang-server/internal/core/transport/ws/middleware"
	presence_redis_repository "github.com/CascadePro/api-golang-server/internal/features/presence/repository/redis"

	"github.com/gorilla/websocket"
)

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
