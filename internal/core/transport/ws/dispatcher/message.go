package core_ws_dispatcher

import (
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type Message struct {
	ID   string                   `json:"id,omitempty"`
	Type domain.RealtimeEventType `json:"type"`
	Data json.RawMessage          `json:"data,omitempty"`
}

func (m *Message) Validate() error {
	if m.ID == "" {
		return fmt.Errorf("`ID` is required: %w", core_errors.ErrInvalidArgument)
	}

	if m.Type == domain.RealtimeEventNil {
		return fmt.Errorf("``Type` is required: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
