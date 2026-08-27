package core_postgres_token

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type TokenModel struct {
	ID      uuid.UUID
	Version int64

	Token string
	Type  domain.TokenType

	ExpiresAt time.Time
	UserID    uuid.UUID
}

func domainFromModel(model TokenModel) domain.Token {
	return domain.Token{
		ID:        model.ID,
		Version:   model.Version,
		Token:     model.Token,
		Type:      model.Type,
		ExpiresAt: model.ExpiresAt,
		UserID:    model.UserID,
	}
}
