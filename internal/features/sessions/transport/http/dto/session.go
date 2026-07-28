package sessions_http_dto

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

type SessionDto struct {
	ID           string             `json:"id"             example:"831f21c1798a7972fa9cda12dac0"`
	IP           string             `json:"ip"             example:"0.0.0.0"`
	CreatedAt    time.Time          `json:"created_at"     example:"2006-01-02T15-04-05.000000"`
	LastActiveAt time.Time          `json:"last_active_at" example:"2006-01-02T15-04-05.000000"`
	ExpiresAt    time.Time          `json:"expires_at"     example:"2006-01-02T15-04-05.000000"`
	Metadata     SessionDtoMetadata `json:"metadata"`
}

type SessionDtoMetadata struct {
	Location SessionDtoMetadataLocation `json:"location"`
	Device   SessionDtoMetadataDevice   `json:"device"`
}

type SessionDtoMetadataLocation struct {
	Country string  `json:"country,omitempty" example:"Russia"`
	City    string  `json:"city,omitempty"    example:"Pskov"`
	Lat     float64 `json:"lat,omitempty"     example:"0.00000"`
	Lng     float64 `json:"lng,omitempty"     example:"0.00000"`
}

type SessionDtoMetadataDevice struct {
	OS      string `json:"os,omitempty"    example:"Windows 10"`
	Model   string `json:"model,omitempty" example:"BOM-WXX9"`
	AppName string `json:"app,omitempty"   example:"Cascade Pro"`
	Type    string `json:"type,omitempty"  example:"desktop"`
	Version string `json:"vers,omitempty"  example:"1.0.0"`
}

func SessionDomainToDTO(session domain.Session) SessionDto {
	return SessionDto{
		ID:           session.ID,
		IP:           session.IP.String(),
		CreatedAt:    session.CreatedAt,
		LastActiveAt: session.LastActiveAt,
		ExpiresAt:    session.CreatedAt.Add(time.Duration(session.ExpirationTime)),
		Metadata: SessionDtoMetadata{
			Location: SessionDtoMetadataLocation(session.Metadata.Location),
			Device:   SessionDtoMetadataDevice(session.Metadata.Device),
		},
	}
}

func SessionDomainsToDTOs(domains []domain.Session) []SessionDto {
	sessions := []SessionDto{}

	for _, session := range domains {
		sessions = append(sessions, SessionDomainToDTO(session))
	}

	return sessions
}
