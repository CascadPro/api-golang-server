package core_validation

import (
	"encoding/hex"
	"fmt"
	"net/mail"
	"strings"
	"unicode"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

const (
	NameMinLen int = 2
	NameMaxLen int = 100

	PasswordMinLen int = 8
	PasswordMaxLen int = 100
)

func ValidateStringLength(v *string, field string, min int, max int) error {
	if v == nil {
		return fmt.Errorf("`%s` can't be NULL: %w", field, core_errors.ErrInvalidArgument)
	}

	if vLen := len([]rune(*v)); vLen < min || vLen > max {
		return fmt.Errorf("invalid `%s` len: %d: %w", field, vLen, core_errors.ErrInvalidArgument)
	}

	return nil
}

func ValidateStringEmail(email *string) (*mail.Address, error) {
	if email == nil {
		return nil, fmt.Errorf("`Email` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}

	address, err := mail.ParseAddress(*email)
	if err != nil {
		return nil, fmt.Errorf("invalid `Email` field: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return address, nil
}

func ValidatePassword(pwd string, fieldName *string) error {
	fn := *fieldName
	if fn == "" {
		fn = "Password"
	}

	pwdLen := len([]rune(pwd))
	if pwdLen < PasswordMinLen || pwdLen > PasswordMaxLen {
		return fmt.Errorf("invalid `%s` len: %d: %w", fn, pwdLen, core_errors.ErrInvalidArgument)
	}

	var hasDigit, hasUpper, hasSpec bool

	for _, r := range pwd {
		switch {
		case !hasDigit && r >= '0' && r <= '9':
			hasDigit = true
		case !hasUpper && unicode.IsUpper(r):
			hasUpper = true
		case !hasSpec && strings.ContainsRune("!@#$%^&*_-", r):
			hasSpec = true
		}
	}

	if !hasDigit {
		return fmt.Errorf("`%s` should contain at least one digit: %w", fn, core_errors.ErrInvalidArgument)
	}

	if !hasUpper {
		return fmt.Errorf("`%s` should contain at least one uppercase letter: %w", fn, core_errors.ErrInvalidArgument)
	}

	if !hasSpec {
		return fmt.Errorf("`%s` should contain at least one special symbol: %w", fn, core_errors.ErrInvalidArgument)
	}

	return nil
}

func ValidateArray[T ~string](list []T, item T) error {
	dict := map[T]struct{}{}
	for _, li := range list {
		dict[li] = struct{}{}
	}

	if _, ok := dict[item]; !ok {
		return fmt.Errorf("item `%v` is out of range: %w", item, core_errors.ErrInvalidArgument)
	}

	return nil
}

func ValidateID(id string, byteLength int) error {
	if id == "" {
		return fmt.Errorf("session id is empty: %w", core_errors.ErrInvalidArgument)
	}

	if byteLength > 0 {
		expected := byteLength * 2
		if len(id) != expected {
			return fmt.Errorf(
				"invalid session id length: got %d, want %d (hex of %d bytes): %w",
				len(id), expected, byteLength, core_errors.ErrInvalidArgument,
			)
		}
	}

	if _, err := hex.DecodeString(id); err != nil {
		return fmt.Errorf("id is not valid hex: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}
