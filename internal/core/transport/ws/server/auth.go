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
	conn.SetReadDeadline(time.Now().Add(authTimeout))

	var message authRequest
	if err := conn.ReadJSON(&message); err != nil {
		return nil, fmt.Errorf("read client message: %w", core_errors.ErrInvalidArgument)
	}

	if message.Type != domain.RealtimeEventAuthRequest || message.Token == "" {
		return nil, core_errors.ErrInvalidArgument
	}

	claims, err := s.authenticator.Authenticate(conn.Context(), message.Token)
	if err != nil {
		return nil, fmt.Errorf("authenticate: %v: %w", err, core_errors.ErrUnauthorized)
	}

	response := authResponse{
		Type:    domain.RealtimeEventAuthSuccess,
		Message: "authenticated",
	}

	conn.SetReadDeadline(time.Time{})
	conn.WriteJSON(response)

	return claims, nil
}
