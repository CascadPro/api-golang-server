package domain

import (
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

type TokenType string

const (
	TokenTypeRegister    = TokenType("register")
	TokenTypeEmailVerify = TokenType("email_verify")
	TokenTypeNil         = TokenType("")
)

var (
	TokenTypes = [2]TokenType{TokenTypeRegister, TokenTypeEmailVerify}
)

type Token struct {
	ID      uuid.UUID
	Version int64

	Token string
	Type  TokenType

	ExpiresAt time.Time
	CreatedAt time.Time

	UserID uuid.UUID
}

func NewToken(token string, tokenType TokenType, expiresAt time.Time, userID uuid.UUID) Token {
	return Token{
		ID:        UninitializedUUID,
		Version:   UninitializedVersion,
		Token:     token,
		Type:      tokenType,
		ExpiresAt: expiresAt,
		UserID:    userID,
	}
}

func NewGetToken(token string, tokenType TokenType, userID uuid.UUID) Token {
	return Token{
		ID:      UninitializedUUID,
		Version: UninitializedVersion,
		Token:   token,
		Type:    tokenType,
		UserID:  userID,
	}
}

func (t *Token) Validate() error {
	if t.CreatedAt.After(t.ExpiresAt) {
		return fmt.Errorf("`ExpiresAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	if t.Type == TokenTypeRegister {
		if _, err := uuid.Parse(t.Token); err != nil {
			return fmt.Errorf("`Token` must be valid UUID if `Type` = '%s': %w",
				TokenTypeRegister, core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
