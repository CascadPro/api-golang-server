package core_ws_client

import (
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_logger "github.com/CascadePro/api-golang-server/internal/core/logger"
	core_ws_response "github.com/CascadePro/api-golang-server/internal/core/transport/ws/response"
	"go.uber.org/zap"
)

func (c *Client) SendEvent(event domain.RealtimeEvent) bool {
	log := core_logger.FromContext(c.conn.Context())

	data, err := marshal(event)
	if err != nil {
		log.Error("marshal event", zap.String("event_id", event.ID), zap.Error(err))

		return false
	}

	return c.enqueue(ChannelMessage{data: data})
}

func (c *Client) Kick(event domain.RealtimeEvent) bool {
	log := core_logger.FromContext(c.conn.Context())

	data, err := marshal(event)
	if err != nil {
		log.Error("marshal event", zap.String("event_id", event.ID), zap.Error(err))

		return false
	}

	return c.enqueue(ChannelMessage{
		data:        data,
		closeAfter:  true,
		closeCode:   CloseCodeSessionRevoked,
		closeReason: CloseReasonSessionRevoked,
	})
}

func (c *Client) enqueue(message ChannelMessage) bool {
	select {
	case <-c.done:
		return false

	case c.send <- message:
		return true
	}
}

func marshal(event domain.RealtimeEvent) ([]byte, error) {
	data, err := json.Marshal(event.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal data: %w", err)
	}

	response := core_ws_response.NewResponse(event.ID, event.Type, data)

	return json.Marshal(response)
}
