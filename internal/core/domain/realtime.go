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

	RealtimeEventPing RealtimeEventType = "realtime.ping"
	RealtimeEventPong RealtimeEventType = "realtime.pong"

	RealtimeEventPresenceOnline  = RealtimeEventType("presence.online")
	RealtimeEventPresenceOffline = RealtimeEventType("presence.offline")

	RealtimeEventSessionRevoked  = RealtimeEventType("session.revoked")
	RealtimeEventSessionsRevoked = RealtimeEventType("session.revoked.all")
	RealtimeEventSessionCreated  = RealtimeEventType("session.created")
	RealtimeEventSessionUpdated  = RealtimeEventType("session.updated")

	RealtimeEventSettingsUpdated = RealtimeEventType("settings.updated")

	RealtimeEventUserUpdated       = RealtimeEventType("user.updated")
	RealtimeEventUserAvatarUpdated = RealtimeEventType("user.avatar.updated")

	RealtimeEventRequestCreated    = RealtimeEventType("request.created")
	RealtimeEventRequestUpdated    = RealtimeEventType("request.updated")
	RealtimeEventRequestStatusChanged = RealtimeEventType("request.status_changed")
	RealtimeEventRequestFileUploaded = RealtimeEventType("request.file_uploaded")
	RealtimeEventRequestFileDeleted  = RealtimeEventType("request.file_deleted")

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
		ID:        uuid.NewString(),
		Type:      eventType,
		UserID:    userID,
		SessionID: sessionID,
		Data:      data,
	}
}

func (e *RealtimeEvent) Validate() error {
	if e.ID == "" {
		return fmt.Errorf("`ID` can't be empty: %w", core_errors.ErrInvalidArgument)
	}

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

type RealtimeEventRequestCreatedData struct {
	ID          string                 `json:"id"`
	Title       string                 `json:"title"`
	Status      RequestStatus          `json:"status"`
	ClientID    *string                `json:"client_id,omitempty"`
	CreatedAt   time.Time              `json:"created_at"`
}

func NewRealtimeEventRequestCreatedData(id string, title string, status RequestStatus, clientID *string, createdAt time.Time) RealtimeEventRequestCreatedData {
	return RealtimeEventRequestCreatedData{
		ID:        id,
		Title:     title,
		Status:    status,
		ClientID:  clientID,
		CreatedAt: createdAt,
	}
}

type RealtimeEventRequestUpdatedData struct {
	ID        string     `json:"id"`
	Title     *string    `json:"title,omitempty"`
	Status    *RequestStatus `json:"status,omitempty"`
	ClientID  *string    `json:"client_id,omitempty"`
	Deadline  *time.Time `json:"deadline,omitempty"`
	WorkTypes []string   `json:"work_types,omitempty"`
	Geography []string   `json:"geography,omitempty"`
	Origin    []RequestOrigin `json:"origin,omitempty"`
	UpdatedAt time.Time  `json:"updated_at"`
}

func NewRealtimeEventRequestUpdatedData(id string, updatedAt time.Time) RealtimeEventRequestUpdatedData {
	return RealtimeEventRequestUpdatedData{
		ID:        id,
		UpdatedAt: updatedAt,
	}
}

type RealtimeEventRequestStatusChangedData struct {
	ID        string       `json:"id"`
	Status    RequestStatus `json:"status"`
	StatusBy  string       `json:"status_by"`
	UpdatedAt time.Time    `json:"updated_at"`
}

func NewRealtimeEventRequestStatusChangedData(id string, status RequestStatus, statusBy string, updatedAt time.Time) RealtimeEventRequestStatusChangedData {
	return RealtimeEventRequestStatusChangedData{
		ID:        id,
		Status:    status,
		StatusBy:  statusBy,
		UpdatedAt: updatedAt,
	}
}

type RealtimeEventRequestFileData struct {
	ID        string    `json:"id"`
	RequestID string    `json:"request_id"`
	Slot      int       `json:"slot"`
	Filename  string    `json:"filename"`
	Size      int64     `json:"size"`
	MimeType  string    `json:"mime_type"`
	CreatedAt time.Time `json:"created_at"`
}

func NewRealtimeEventRequestFileData(id, requestID string, slot int, filename string, size int64, mimeType string, createdAt time.Time) RealtimeEventRequestFileData {
	return RealtimeEventRequestFileData{
		ID:        id,
		RequestID: requestID,
		Slot:      slot,
		Filename:  filename,
		Size:      size,
		MimeType:  mimeType,
		CreatedAt: createdAt,
	}
}
