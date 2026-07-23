package core_http_request

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
)

type UserAgent struct {
	AppName    string
	AppVersion string

	OS    string
	Model string
	Type  string
}

func ParseUserAgent(r *http.Request) (*UserAgent, error) {
	ua := strings.TrimSpace(r.UserAgent())

	re := regexp.MustCompile(`^([^/]+)/([^(]+)\s*\(([^;]+);\s*([^;]+);\s*([^)]+)\)$`)

	matches := re.FindStringSubmatch(ua)
	if matches == nil {
		return nil, fmt.Errorf("invalid user-agent header format: %w", core_errors.ErrInvalidArgument)
	}

	deviceType := strings.TrimSpace(matches[5])

	switch deviceType {
	case "1":
		deviceType = "phone"
	case "2":
		deviceType = "tablet"
	case "3":
		deviceType = "desktop"
	default:
		deviceType = "unknown"
	}

	return &UserAgent{
		AppName:    strings.TrimSpace(matches[1]),
		AppVersion: strings.TrimSpace(matches[2]),
		OS:         strings.TrimSpace(matches[3]),
		Model:      strings.TrimSpace(matches[4]),
		Type:       deviceType,
	}, nil
}
