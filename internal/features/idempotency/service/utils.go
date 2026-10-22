package idempotency_service

import (
	"net/http"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
)

var (
	headers = []string{"Content-Type", "Content-Language", "Location", "ETag", "Cache-Control"}
)

func equalBytes(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}

	var diff byte
	for i := range a {
		diff |= a[i] ^ b[i]
	}

	return diff == 0
}

func safeHeaders(header http.Header) http.Header {
	result := make(http.Header)

	for _, name := range headers {
		if values := header.Values(name); len(values) > 0 {
			result[name] = append([]string(nil), values...)
		}
	}

	return result
}

func copyHeaders(dst, src http.Header) {
	for name, values := range src {
		dst.Del(name)

		for _, value := range values {
			dst.Add(name, value)
		}
	}
}

func replayResponse(w http.ResponseWriter, record domain.IdempotencyRecord) {
	copyHeaders(w.Header(), record.ResponseHeaders)
	w.Header().Set(HeaderReplayed, "true")

	status := record.ResponseStatus
	if status == 0 {
		status = http.StatusOK
	}

	w.WriteHeader(status)

	if status == http.StatusNoContent || status == http.StatusNotModified {
		return
	}

	_, _ = w.Write(record.ResponseBody)
}
