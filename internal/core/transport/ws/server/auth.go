package core_ws_server

import (
	"fmt"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_jwt_security "github.com/CascadePro/api-golang-server/internal/core/security/jwt"
	core_ws_conn "github.com/CascadePro/api-golang-server/internal/core/transport/ws/conn"
)

const authTimeout = 5 * time.Second

type authRequest struct {
	Type  domain.RealtimeEventType `json:"type"`
	Token string                   `json:"token"`
}

type authResponse struct {
	Type    domain.RealtimeEventType `json:"type"`
	Message string                   `json:"message"`
}

func (s *Server) AuthorizeClient(conn *core_ws_conn.Conn) (*core_jwt_security.AccessClaims, error) {
	if err := conn.SetReadDeadline(time.Now().Add(authTimeout)); err != nil {
		return nil, err
	}

	var message authRequest

	if err := conn.ReadJSON(&message); err != nil {
		return nil, fmt.Errorf("read authentication message: %w", core_errors.ErrInvalidArgument)
	}

	if message.Type != domain.RealtimeEventAuthRequest {
		return nil, fmt.Errorf("invalid authentication event: %w", core_errors.ErrInvalidArgument)
	}

	if message.Token == "" {
		return nil, fmt.Errorf("authentication token is required: %w", core_errors.ErrUnauthorized)
	}

	claims, err := s.authenticator.Authenticate(conn.Context(), message.Token)
	if err != nil {
		return nil, fmt.Errorf("authenticate websocket client: %w", err)
	}

	if err := conn.SetReadDeadline(time.Time{}); err != nil {
		return nil, err
	}

	response := authResponse{
		Type:    domain.RealtimeEventAuthSuccess,
		Message: "authenticated",
	}

	if err := conn.WriteJSON(response); err != nil {
		return nil, fmt.Errorf("write authentication response: %w", err)
	}

	return claims, nil
}
