package core_http_utils

import (
	"fmt"
	"net"
	"net/http"
	"slices"
	"strings"
)

func ParseIP(value string) net.IP {
	value = strings.TrimSpace(value)
	if host, _, err := net.SplitHostPort(value); err == nil {
		value = host
	}

	return net.ParseIP(value)
}

// ClientIPFromRequest trusts forwarding headers only when the direct peer is trusted.
// trustedProxyCIDRs accepts either individual IPs or CIDR ranges.
func ClientIPFromRequest(r *http.Request, trustedProxyCIDRs []string) (net.IP, error) {
	remoteIP := ParseIP(r.RemoteAddr)
	if remoteIP == nil {
		return nil, fmt.Errorf("parse remote address %q", r.RemoteAddr)
	}

	trusted, err := parseTrustedProxies(trustedProxyCIDRs)
	if err != nil {
		return nil, err
	}

	fmt.Println(trusted)
	fmt.Println(remoteIP.String())

	if !containsIP(trusted, remoteIP) {
		return remoteIP, nil
	}

	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		parts := strings.Split(forwarded, ",")
		chain := make([]net.IP, 0, len(parts)+1)

		for _, part := range parts {
			ip := ParseIP(part)
			if ip == nil {
				return nil, fmt.Errorf("invalid IP in X-Forwarded-For")
			}

			chain = append(chain, ip)
		}

		for _, c := range slices.Backward(chain) {
			if !containsIP(trusted, c) {
				return c, nil
			}
		}

		return nil, fmt.Errorf("X-Forwarded-For contains only trusted proxies")
	}

	if realIP := r.Header.Get("X-Real-IP"); realIP != "" {
		ip := ParseIP(realIP)
		if ip == nil {
			return nil, fmt.Errorf("invalid X-Real-IP")
		}

		return ip, nil
	}

	return remoteIP, nil
}

// ClientIP is retained for callers not yet migrated. Without an explicit trust list,
// forwarded headers are deliberately ignored.
func ClientIP(r *http.Request) (net.IP, error) {
	return ClientIPFromRequest(r, nil)
}

func parseTrustedProxies(values []string) ([]*net.IPNet, error) {
	out := make([]*net.IPNet, 0, len(values))
	for _, value := range values {
		value = strings.TrimSpace(value)
		if value == "" {
			continue
		}

		if ip := net.ParseIP(value); ip != nil {
			bits := 128
			if ip.To4() != nil {
				ip = ip.To4()
				bits = 32
			}

			mask := net.CIDRMask(bits, bits)

			out = append(out, &net.IPNet{IP: ip, Mask: mask})

			continue
		}

		_, network, err := net.ParseCIDR(value)
		if err != nil {
			return nil, fmt.Errorf("invalid trusted proxy %q: %w", value, err)
		}

		out = append(out, network)
	}

	return out, nil
}

func containsIP(networks []*net.IPNet, ip net.IP) bool {
	for _, network := range networks {
		if network.Contains(ip) {
			return true
		}
	}

	return false
}
