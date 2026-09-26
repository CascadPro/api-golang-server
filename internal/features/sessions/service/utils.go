package session_service

import (
	"encoding/json"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

func marshalOutboxEventPayload(
	t domain.RealtimeEventType,
	userID *uuid.UUID,
	sessionID string,
) ([]byte, error) {
	payload, err := domain.NewEventRealtimePublishPayload(
		domain.RealtimeEventSessionsRevoked,
		userID,
		nil,
		domain.NewRealtimeEventSessionRevokeData(
			sessionID,
			domain.SessionRevokeDataUserRevoked,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("new outbox event: %w", err)
	}

	result, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshal payload: %w", err)
	}

	return result, nil
}
