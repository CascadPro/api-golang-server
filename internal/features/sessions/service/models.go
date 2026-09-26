package session_service

import "github.com/CascadePro/api-golang-server/internal/core/domain"

type Session struct {
	domain.Session

	Online bool
}
