package core_ws_dispatcher

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type Message struct {
	ID        string                   `json:"id,omitempty"`
	Type      domain.RealtimeEventType `json:"type"`
	Timestamp *time.Time               `json:"timestamp,omitempty"`
	Data      json.RawMessage          `json:"data,omitempty"`
}

func (m *Message) Validate() error {
	if m.Type == domain.RealtimeEventNil {
		return fmt.Errorf("``Type` is required: %w", core_errors.ErrInvalidArgument)
	}

	if m.Timestamp != nil && time.Now().Before(*m.Timestamp) {
		return fmt.Errorf("`Timestamp` can't be more than current time: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
