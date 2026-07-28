package core_postgres_token

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type TokenModel struct {
	ID      uuid.UUID
	Version int

	Token string
	Type  domain.TokenType

	ExpiresAt time.Time
	UserID    uuid.UUID
}
