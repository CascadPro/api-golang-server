package domain

import (
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

type RealtimeEventType string

const (
	RealtimeEventAuthRequest = RealtimeEventType("auth.request")
	RealtimeEventAuthSuccess = RealtimeEventType("auth.success")
	RealtimeEventAuthError   = RealtimeEventType("auth.error")

	RealtimeEventPresenceOnline  = RealtimeEventType("presence.online")
	RealtimeEventPresenceOffline = RealtimeEventType("presence.offline")

	RealtimeEventSessionRevoked  = RealtimeEventType("session.revoked")
	RealtimeEventSessionsRevoked = RealtimeEventType("session.revoked.all")
	RealtimeEventSessionCreated  = RealtimeEventType("session.created")
	RealtimeEventSessionUpdated  = RealtimeEventType("session.updated")

	RealtimeEventNil = RealtimeEventType("")
)

type RealtimeEvent struct {
	ID         string
	Type       RealtimeEventType
	UserID     *uuid.UUID
	SessionID  *string
	Recipients []uuid.UUID
	Data       any
	CreatedAt  time.Time
}

func NewRealtimeEvent(eventType RealtimeEventType, userID *uuid.UUID, sessionID *string, data any) RealtimeEvent {
	return RealtimeEvent{
		ID:        UninitializedID,
		Type:      eventType,
		UserID:    userID,
		SessionID: sessionID,
		Data:      data,
	}
}

func (e *RealtimeEvent) Validate() error {
	if e.Type == RealtimeEventNil {
		return fmt.Errorf("`Type` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if e.UserID != nil && *e.UserID == UninitializedUUID {
		return fmt.Errorf("`UserID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	if e.SessionID != nil && *e.SessionID == UninitializedID {
		return fmt.Errorf("`SessionID` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type RealtimeEventPresenceData struct {
	UserID    uuid.UUID `json:"uid"`
	SessionID string    `json:"sid"`
}

func NewRealtimeEventPresenceData(userID uuid.UUID, sessionID string) RealtimeEventPresenceData {
	return RealtimeEventPresenceData{
		UserID:    userID,
		SessionID: sessionID,
	}
}

type SessionRevokeDataReason string

const (
	SessionRevokeDataUserRevoked = SessionRevokeDataReason("user_revoked")
)

type RealtimeEventSessionRevokeData struct {
	SessionID string                  `json:"sid"`
	Reason    SessionRevokeDataReason `json:"reason"`
}

func NewRealtimeEventSessionRevokeData(sessionID string, reason SessionRevokeDataReason) RealtimeEventSessionRevokeData {
	return RealtimeEventSessionRevokeData{
		SessionID: sessionID,
		Reason:    reason,
	}
}
