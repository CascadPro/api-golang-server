package core_ws_conn

import (
	"context"

	"github.com/gorilla/websocket"
)

type Conn struct {
	*websocket.Conn

	ctx context.Context
}

func NewConn(ctx context.Context, conn *websocket.Conn) *Conn {
	return &Conn{
		Conn: conn,
		ctx:  ctx,
	}
}

func (c *Conn) Context() context.Context {
	return c.ctx
}

func (c *Conn) WithContext(ctx context.Context) *Conn {
	c.ctx = ctx
	return c
}
