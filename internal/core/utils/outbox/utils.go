package core_outbox_utils

import (
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

func marshalPayload(
	t domain.RealtimeEventType,
	userID *uuid.UUID,
	extra any,
) (json.RawMessage, error) {
	payload, err := domain.NewEventRealtimePublishPayload(t, userID, nil, extra)
	if err != nil {
		return nil, fmt.Errorf("new outbox event: %w", err)
	}

	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	return data, nil
}
