package core_http_utils

import (
	"net/http/httptest"
	"testing"
)

func TestClientIPFromRequestIgnoresUntrustedForwardedHeaders(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "203.0.113.9:1234"
	r.Header.Set("X-Forwarded-For", "198.51.100.7")

	ip, err := ClientIPFromRequest(r, []string{"10.0.0.0/8"})
	if err != nil {
		t.Fatal(err)
	}
	if got := ip.String(); got != "203.0.113.9" {
		t.Fatalf("expected remote IP, got %s", got)
	}
}

func TestClientIPFromRequestUsesForwardedChainOnlyFromTrustedProxy(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "10.0.0.2:8080"
	r.Header.Set("X-Forwarded-For", "198.51.100.7, 10.0.0.1")

	ip, err := ClientIPFromRequest(r, []string{"10.0.0.0/8"})
	if err != nil {
		t.Fatal(err)
	}
	if got := ip.String(); got != "198.51.100.7" {
		t.Fatalf("expected client IP, got %s", got)
	}
}

func TestClientIPFromRequestRejectsMalformedForwardedIPFromTrustedProxy(t *testing.T) {
	r := httptest.NewRequest("GET", "/", nil)
	r.RemoteAddr = "10.0.0.2:8080"
	r.Header.Set("X-Forwarded-For", "not-an-ip")

	if _, err := ClientIPFromRequest(r, []string{"10.0.0.0/8"}); err == nil {
		t.Fatal("expected malformed forwarded IP error")
	}
}
