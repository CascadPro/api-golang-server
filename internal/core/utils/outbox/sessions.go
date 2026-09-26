package core_outbox_utils

import (
	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

func NewRealtimeEventSessionRevoked(
	userID uuid.UUID,
	sessionID string,
) (domain.OutboxEvent, error) {
	return newSessionRevoked(domain.RealtimeEventSessionRevoked, &userID, &sessionID)
}

func NewRealtimeEventSessionsRevoked(
	userID uuid.UUID,
	sessionID string,
) (domain.OutboxEvent, error) {
	return newSessionRevoked(domain.RealtimeEventSessionsRevoked, &userID, &sessionID)
}

func newSessionRevoked(
	t domain.RealtimeEventType,
	userID *uuid.UUID,
	sessionID *string,
) (domain.OutboxEvent, error) {
	data := domain.NewRealtimeEventSessionRevokeData(
		*sessionID,
		domain.SessionRevokeDataUserRevoked,
	)

	payload, err := marshalPayload(t, userID, data)
	if err != nil {
		return domain.OutboxEvent{}, err
	}

	event := domain.NewOutboxEvent(
		domain.EventTypeRealtimePublish,
		userID,
		payload,
	)

	return event, nil
}
