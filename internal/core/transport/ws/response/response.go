package core_ws_response

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

type Response struct {
	ID        string                   `json:"id"`
	Type      domain.RealtimeEventType `json:"type"`
	Data      []byte                   `json:"data,omitempty"`
	Timestamp time.Time                `json:"timestamp"`
}

func NewResponse(id string, t domain.RealtimeEventType, data []byte) *Response {
	return &Response{
		ID:        id,
		Type:      t,
		Data:      data,
		Timestamp: time.Now().UTC(),
	}
}
