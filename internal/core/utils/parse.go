package core_utils

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

// ParseDurationExtended can parse strings:
//
//	30d  → 30 * 24h
//	2w   → 2  * 7d
//	5h   → hours (like a time.ParseDuration)
//	10m  → minutes
//	15s  → seconds and other
func ParseDurationExtended(s string) (time.Duration, error) {
	if d, err := time.ParseDuration(s); err == nil {
		return d, nil
	}

	s = strings.TrimSpace(strings.ToLower(s))

	mult := map[string]time.Duration{
		"ns": time.Nanosecond,
		"us": time.Microsecond,
		"µs": time.Microsecond,
		"ms": time.Millisecond,
		"s":  time.Second,
		"m":  time.Minute,
		"h":  time.Hour,
		"d":  24 * time.Hour,
		"w":  7 * 24 * time.Hour,
	}

	var numPart, unitPart string
	for i, r := range s {
		if r < '0' || r > '9' {
			numPart = s[:i]
			unitPart = s[i:]
			break
		}
	}
	if numPart == "" || unitPart == "" {
		return 0, fmt.Errorf("invalid duration %q: %w", s, core_errors.ErrInvalidArgument)
	}

	val, err := strconv.ParseInt(numPart, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid number %q in %q: %v: %w", numPart, s, err, core_errors.ErrInvalidArgument)
	}

	m, ok := mult[unitPart]
	if !ok {
		return 0, fmt.Errorf("unsupported duration unit '%q': %w", unitPart, core_errors.ErrInvalidArgument)
	}

	return time.Duration(val) * m, nil
}

func ParseFloat(s string) float64 {
	var f float64

	if _, err := fmt.Sscanf(s, "%f", &f); err != nil {
		return 0
	}

	return f
}
