package domain

import (
	"errors"
	"fmt"
	"strings"

	pn "github.com/nyaruka/phonenumbers"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type PhoneNumber struct {
	E164         string `json:"e164"`
	CountryCode  string `json:"country_code"`
	OperatorCode string `json:"operator_code"`
	Subscriber   string `json:"subscriber"`
}

func (p PhoneNumber) IsZero() bool {
	return p.E164 == ""
}

func (p PhoneNumber) String() string {
	return p.E164
}

func (p PhoneNumber) MarshalJSON() ([]byte, error) {
	if p.IsZero() {
		return []byte("null"), nil
	}
	return []byte(`"` + p.E164 + `"`), nil
}

func (p *PhoneNumber) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	if s == "null" || s == "" {
		*p = PhoneNumber{}
		return nil
	}

	parsed, err := ParsePhoneNumber(s)
	if err != nil {
		return err
	}

	*p = parsed
	return nil
}

// ParsePhoneNumber parses a phone number string and returns a structured PhoneNumber.
// Accepts various formats: E.164 (+79001234567), national (89001234567), with spaces, etc.
// Default region is empty string (auto-detect from + prefix).
func ParsePhoneNumber(raw string) (PhoneNumber, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return PhoneNumber{}, fmt.Errorf("phone number is empty: %w", core_errors.ErrInvalidArgument)
	}

	num, err := pn.Parse(raw, "")
	if err != nil {
		return PhoneNumber{}, fmt.Errorf("cannot parse phone number: %v: %w", err, core_errors.ErrInvalidArgument)
	}

	if !pn.IsValidNumber(num) {
		return PhoneNumber{}, fmt.Errorf("phone number is not valid: %w", core_errors.ErrInvalidArgument)
	}

	e164 := pn.Format(num, pn.E164)

	countryCode := fmt.Sprintf("%d", num.GetCountryCode())
	national := fmt.Sprintf("%d", num.GetNationalNumber())

	// Operator code lengths by country
	// Russia (7): 3 digits (e.g., 900, 916, 999)
	// Belarus (375): 3 digits (e.g., 29, 33, 44 - but these are 2 digits, handled below)
	operatorPrefixLen := map[string]int{
		"7":   3, // Russia
		"375": 2, // Belarus (operator codes are 2 digits: 29, 33, 44, 25)
		"1":   3, // US/Canada (area code)
		"44":  4, // UK (area codes vary, typically 3-5 digits after leading 0 removed)
	}

	prefixLen := 0
	if l, ok := operatorPrefixLen[countryCode]; ok {
		prefixLen = l
	} else {
		if len(national) > 7 {
			prefixLen = 3
		}
	}

	var operatorCode, subscriber string
	if prefixLen > 0 && len(national) > prefixLen {
		operatorCode = national[:prefixLen]
		subscriber = national[prefixLen:]
	} else {
		subscriber = national
	}

	return PhoneNumber{
		E164:         e164,
		CountryCode:  countryCode,
		OperatorCode: operatorCode,
		Subscriber:   subscriber,
	}, nil
}

func MustParsePhoneNumber(raw string) PhoneNumber {
	p, err := ParsePhoneNumber(raw)
	if err != nil {
		panic(err)
	}

	return p
}

func (p PhoneNumber) FormatNational(region string) (string, error) {
	if p.IsZero() {
		return "", errors.New("phone number is zero")
	}

	num, err := pn.Parse(p.E164, region)
	if err != nil {
		return "", err
	}

	return pn.Format(num, pn.NATIONAL), nil
}

func (p PhoneNumber) FormatInternational() (string, error) {
	if p.IsZero() {
		return "", errors.New("phone number is zero")
	}

	num, err := pn.Parse(p.E164, "")
	if err != nil {
		return "", err
	}

	return pn.Format(num, pn.INTERNATIONAL), nil
}

func ValidatePhoneNumber(phone *string) (PhoneNumber, error) {
	if phone == nil {
		return PhoneNumber{}, fmt.Errorf("phone number cannot be nil: %w", core_errors.ErrInvalidArgument)
	}
	return ParsePhoneNumber(*phone)
}

func ValidatePhoneNumberNullable(phone Nullable[string]) (PhoneNumber, bool, error) {
	if !phone.Set {
		return PhoneNumber{}, false, nil
	}

	if phone.Value == nil {
		return PhoneNumber{}, true, fmt.Errorf("phone number cannot be null when set: %w", core_errors.ErrInvalidArgument)
	}

	p, err := ParsePhoneNumber(*phone.Value)
	return p, true, err
}

func PhoneNumberFromE164(e164 string) (PhoneNumber, error) {
	return ParsePhoneNumber(e164)
}
