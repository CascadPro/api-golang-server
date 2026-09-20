package core_ws_client

import (
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type ChannelMessage struct {
	data        []byte
	closeAfter  bool
	closeCode   CloseCode
	closeReason CloseReason
}

type ClientMessage struct {
	ID        string                   `json:"id,omitempty"`
	Type      domain.RealtimeEventType `json:"type"`
	Timestamp string                   `json:"timestamp,omitempty"`
	Data      json.RawMessage          `json:"data,omitempty"`
}

func (m *ClientMessage) Validate() error {
	if m.Type == domain.RealtimeEventNil {
		return fmt.Errorf("`type` is required: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
