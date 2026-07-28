package core_http_utils

import (
	"fmt"
	"net"
	"net/http"
	"strings"

	core_config "github.com/CascadePro/api-golang-server/internal/core/config"
)

var PskovIP net.IP = ParseIP("109.207.190.162")

// ParseIP проверяет, что строка действительно IP (IPv4 или IPv6)
func ParseIP(s string) net.IP {
	s = strings.TrimSpace(s)

	if strings.Contains(s, ":") && strings.Count(s, ":") == 1 {
		host, _, _ := net.SplitHostPort(s)
		s = host
	}

	ip := net.ParseIP(s)

	return ip
}

// ClientIP возвращает «наиболее вероятный» публичный IP клиента.
// Порядок приоритета:
//  1. X-Forwarded-For (первый в списке)
//  2. X-Real-IP
//  3. RemoteAddr
func ClientIP(r *http.Request) (net.IP, error) {
	cfg, err := core_config.NewConfig()
	if err != nil {
		return nil, fmt.Errorf("get core config: %w", err)
	}

	if cfg.Connection == core_config.ConnectionOffline || cfg.EnvMode == core_config.EnvModeDev {
		return PskovIP, nil
	}

	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			candidate := ParseIP(parts[0])
			if candidate != nil {
				return candidate, nil
			}
		}
	}

	if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
		if ip := ParseIP(xrip); ip != nil {
			return ip, nil
		}
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	if ip := net.ParseIP(host); ip != nil {
		return ip, nil
	}

	return nil, fmt.Errorf("cannot determine client IP")
}
