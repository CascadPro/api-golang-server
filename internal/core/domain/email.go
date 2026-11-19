package domain

import (
	"fmt"
	"net/mail"
	"strings"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type Email string

func ParseEmail(raw string) (Email, error) {
	addr, err := mail.ParseAddress(strings.TrimSpace(raw))
	if err != nil {
		return "", fmt.Errorf("invalid email format: %w", core_errors.ErrInvalidArgument)
	}

	normalized := strings.ToLower(addr.Address)

	return Email(normalized), nil
}

func MustParseEmail(raw string) Email {
	e, err := ParseEmail(raw)
	if err != nil {
		panic(err)
	}

	return e
}

func (e Email) String() string {
	return string(e)
}

func (e Email) IsZero() bool {
	return e == ""
}

func ValidateEmail(email *string) (Email, error) {
	if email == nil {
		return "", fmt.Errorf("email cannot be nil: %w", core_errors.ErrInvalidArgument)
	}
	return ParseEmail(*email)
}

// ValidateEmailNullable validates a Nullable[string] as an email.
// Returns the parsed Email and whether it was set.
func ValidateEmailNullable(email Nullable[string]) (Email, bool, error) {
	if !email.Set {
		return "", false, nil
	}
	if email.Value == nil {
		return "", true, fmt.Errorf("email cannot be null when set: %w", core_errors.ErrInvalidArgument)
	}
	e, err := ParseEmail(*email.Value)
	return e, true, err
}
