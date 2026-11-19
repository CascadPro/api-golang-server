package core_validation

import (
	"encoding/hex"
	"fmt"
	"strings"
	"unicode"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type Validatable interface {
	Validate() error
}

const (
	NameMinLen int = 2
	NameMaxLen int = 100

	PasswordMinLen int = 8
	PasswordMaxLen int = 100
)

func ValidateStringLength(v *string, field string, minLen, maxLen int) error {
	if v == nil {
		return fmt.Errorf("`%s` can't be NULL: %w", field, core_errors.ErrInvalidArgument)
	}

	if vLen := len([]rune(*v)); vLen < minLen || vLen > maxLen {
		return fmt.Errorf("invalid `%s` len: %d: %w", field, vLen, core_errors.ErrInvalidArgument)
	}

	return nil
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
		return fmt.Errorf("id is empty: %w", core_errors.ErrInvalidArgument)
	}

	if byteLength > 0 {
		expected := byteLength * 2
		if len(id) != expected {
			return fmt.Errorf(
				"invalid id length: got %d, want %d (hex of %d bytes): %w",
				len(id), expected, byteLength, core_errors.ErrInvalidArgument,
			)
		}
	}

	if _, err := hex.DecodeString(id); err != nil {
		return fmt.Errorf("id is not valid hex: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	return nil
}

func ValidateInteger(x int, fieldName string, minLen *int, maxLen *int) error {
	if minLen != nil && x < *minLen {
		return fmt.Errorf("`%s` must be greater than or equal %d: %w", fieldName, *minLen, core_errors.ErrInvalidArgument)
	}

	if maxLen != nil && x > *maxLen {
		return fmt.Errorf("`%s` must be less than or equal %d: %w", fieldName, *maxLen, core_errors.ErrInvalidArgument)
	}

	return nil
}
