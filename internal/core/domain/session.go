package domain

import (
	"fmt"
	"net"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
)

type SessionLifetime time.Duration

const (
	SessionLifetime14d = SessionLifetime(time.Hour * 24 * 14)
	SessionLifetime30d = SessionLifetime(time.Hour * 24 * 30)
	SessionLifetime90d = SessionLifetime(SessionLifetime30d * 3)
)

const SessionIDByteLength = 12

type Session struct {
	ID             string
	CreatedAt      time.Time
	LastActiveAt   time.Time
	ExpirationTime SessionLifetime
	IP             net.IP
	Metadata       SessionMetadata
}

type SessionMetadata struct {
	Location SessionMetadataLocation
	Device   SessionMetadataDevice
}

type SessionMetadataLocation struct {
	Country string
	City    string
	Lat     float64
	Lng     float64
}

type SessionMetadataDevice struct {
	OS      string
	Model   string
	AppName string
	Type    string
	Version string
}

func NewSession(id string, ip net.IP, metadata SessionMetadata, createdAt, lastActiveAt, expiresAt time.Time) Session {
	return Session{
		ID:             id,
		IP:             ip,
		Metadata:       metadata,
		CreatedAt:      createdAt,
		LastActiveAt:   lastActiveAt,
		ExpirationTime: SessionLifetime(expiresAt.Sub(createdAt)),
	}
}

func NewAuthSession(ip net.IP, metadata SessionMetadata, expirationTime SessionLifetime) Session {
	return Session{
		ID:             UninitializedID,
		IP:             ip,
		Metadata:       metadata,
		CreatedAt:      time.Now(),
		LastActiveAt:   time.Now(),
		ExpirationTime: expirationTime,
	}
}

func NewSessionMetadata(location SessionMetadataLocation, device SessionMetadataDevice) SessionMetadata {
	return SessionMetadata{
		Location: location,
		Device:   device,
	}
}

func NewSessionMetadataLocation(country, city string, lat, lng float64) SessionMetadataLocation {
	return SessionMetadataLocation{
		Country: country,
		City:    city,
		Lat:     lat,
		Lng:     lng,
	}
}

func NewSessionMetadataDevice(os, model, appName, Type, version string) SessionMetadataDevice {
	return SessionMetadataDevice{
		OS:      os,
		Model:   model,
		AppName: appName,
		Type:    Type,
		Version: version,
	}
}

func (s *Session) Validate() error {
	if len(s.IP) == 0 {
		return fmt.Errorf("`IP` can't be nil or empty: %w", core_errors.ErrInvalidArgument)
	}

	if s.CreatedAt.After(s.CreatedAt.Add(time.Duration(s.ExpirationTime))) {
		return fmt.Errorf("`CreatedAt` can't be after `ExpiresAt`: %w", core_errors.ErrInvalidArgument)
	}

	if err := s.Metadata.Validate(); err != nil {
		return fmt.Errorf("validate metadata: %w", err)
	}

	return nil
}

func (m *SessionMetadata) Validate() error {
	if err := m.Location.Validate(); err != nil {
		return fmt.Errorf("validate location: %w", err)
	}
	if err := m.Device.Validate(); err != nil {
		return fmt.Errorf("validate device: %w", err)
	}

	return nil
}

func (l *SessionMetadataLocation) Validate() error {
	min, max := core_validation.NameMinLen, core_validation.NameMaxLen

	if err := core_validation.ValidateStringLength(&l.Country, "Country", min, max); err != nil {
		return fmt.Errorf("country validation: %w", err)
	}
	if err := core_validation.ValidateStringLength(&l.City, "City", min, max); err != nil {
		return fmt.Errorf("city validation: %w", err)
	}
	if l.Lat < -90.0 || l.Lat > 90.0 {
		return fmt.Errorf("`Lat` must be in range [-90, 90]: %w", core_errors.ErrInvalidArgument)
	}
	if l.Lng < -180.0 || l.Lng > 180.0 {
		return fmt.Errorf("`Lng` must be in range [-180, 180]: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func (d *SessionMetadataDevice) Validate() error {
	min, max := core_validation.NameMinLen, core_validation.NameMaxLen

	if err := core_validation.ValidateStringLength(&d.OS, "OS", min, max); err != nil {
		return fmt.Errorf("os validation: %w", err)
	}
	if err := core_validation.ValidateStringLength(&d.Model, "Model", min, max); err != nil {
		return fmt.Errorf("model validation: %w", err)
	}
	if err := core_validation.ValidateStringLength(&d.AppName, "AppName", min, max); err != nil {
		return fmt.Errorf("app name validation: %w", err)
	}
	if err := core_validation.ValidateStringLength(&d.Type, "Type", min, max); err != nil {
		return fmt.Errorf("type validation: %w", err)
	}
	if err := core_validation.ValidateStringLength(&d.Version, "Version", min, max); err != nil {
		return fmt.Errorf("version validation: %w", err)
	}

	return nil
}
