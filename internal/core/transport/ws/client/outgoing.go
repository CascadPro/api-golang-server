package core_ws_client

import (
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_ws_response "github.com/CascadePro/api-golang-server/internal/core/transport/ws/response"
)

func (c *Client) SendEvent(event domain.RealtimeEvent) bool {
	data, err := marshal(event)
	if err != nil {
		c.logMarshalError(event.ID, event.Type, err)
		return false
	}

	return c.enqueue(ChannelMessage{
		data: data,
	})
}

func (c *Client) RevokeSession(event domain.RealtimeEvent) bool {
	data, err := marshal(event)
	if err != nil {
		c.logMarshalError(event.ID, event.Type, err)
		return false
	}

	return c.enqueue(ChannelMessage{data: data})
}

func (c *Client) RevokeAllSessions(event domain.RealtimeEvent) bool {
	data, err := marshal(event)
	if err != nil {
		c.logMarshalError(event.ID, event.Type, err)
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

	default:
		return false
	}
}

func marshal(event domain.RealtimeEvent) ([]byte, error) {
	if err := event.Validate(); err != nil {
		return nil, fmt.Errorf("validate realtime event: %w", err)
	}

	data, err := json.Marshal(event.Data)
	if err != nil {
		return nil, fmt.Errorf("marshal event data: %w", err)
	}

	response := core_ws_response.NewResponse(
		event.ID,
		event.Type,
		json.RawMessage(data),
	)

	return json.Marshal(response)
}
