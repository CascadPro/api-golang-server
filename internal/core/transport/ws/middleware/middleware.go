package core_ws_middleware

import (
	"slices"

	core_ws_conn "github.com/CascadePro/api-golang-server/internal/core/transport/ws/conn"
)

type Handler func(conn *core_ws_conn.Conn) error

type Middleware func(Handler) Handler

func ChainMiddleware(h Handler, m ...Middleware) Handler {
	if len(m) == 0 {
		return h
	}

	for _, v := range slices.Backward(m) {
		h = v(h)
	}

	return h
}
