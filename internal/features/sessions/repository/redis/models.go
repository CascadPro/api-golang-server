package sessions_redis_repository

import (
	"encoding/json"
	"fmt"
	"net"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_redis_pool "github.com/CascadePro/api-golang-server/internal/core/repository/redis/pool"
)

type SessionModel struct {
	CreatedAt    time.Time `redis:"created_at"`
	ExpiresAt    time.Time `redis:"exp"`
	LastActiveAt time.Time `redis:"last_active_at"`
	IP           net.IP    `redis:"ip"`
	Metadata     []byte    `redis:"metadata"`
}

type SessionModelMetadata struct {
	Location SessionModelMetadataLocation `json:"location"`
	Device   SessionModelMetadataDevice   `json:"device"`
}

type SessionModelMetadataLocation struct {
	Country string  `json:"country"`
	City    string  `json:"city"`
	Lat     float64 `json:"lat"`
	Lng     float64 `json:"lng"`
}

type SessionModelMetadataDevice struct {
	OS      string `json:"os"`
	Model   string `json:"model"`
	AppName string `json:"app"`
	Type    string `json:"type"`
	Version string `json:"vers"`
}

func domainToModel(session domain.Session) (SessionModel, error) {
	metadata := SessionModelMetadata{
		Location: SessionModelMetadataLocation(session.Metadata.Location),
		Device:   SessionModelMetadataDevice(session.Metadata.Device),
	}

	jsonMetadata, err := json.Marshal(metadata)
	if err != nil {
		return SessionModel{}, fmt.Errorf("session model metadata json marshal: %w", err)
	}

	return SessionModel{
		IP:           session.IP,
		CreatedAt:    session.CreatedAt,
		LastActiveAt: session.LastActiveAt,
		ExpiresAt:    session.CreatedAt.Add(time.Duration(session.ExpirationTime)),
		Metadata:     jsonMetadata,
	}, nil
}

func modelToDomain(sessionID string, model SessionModel) (domain.Session, error) {
	if len(model.Metadata) == 0 {
		return domain.Session{}, core_redis_pool.ErrNoValue
	}

	var modelMetadata SessionModelMetadata
	if err := json.Unmarshal(model.Metadata, &modelMetadata); err != nil {
		return domain.Session{}, fmt.Errorf("session model metadata json unmarshal: %w", err)
	}

	metadata := domain.NewSessionMetadata(
		domain.SessionMetadataLocation(modelMetadata.Location),
		domain.SessionMetadataDevice(modelMetadata.Device),
	)

	return domain.NewSession(
		sessionID,
		model.IP,
		metadata,
		model.CreatedAt,
		model.LastActiveAt,
		model.ExpiresAt,
	), nil
}
