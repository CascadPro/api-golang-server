package core_ws_conn

import (
	"context"

	"github.com/gorilla/websocket"
)

type Conn struct {
	*websocket.Conn

	ctx    context.Context
	cancel context.CancelFunc
}

func NewConn(ctx context.Context, cancel context.CancelFunc, conn *websocket.Conn) *Conn {
	return &Conn{
		Conn: conn,

		ctx:    ctx,
		cancel: cancel,
	}
}

func (c *Conn) Context() context.Context {
	return c.ctx
}

func (c *Conn) CancelContext() {
	c.cancel()
}
