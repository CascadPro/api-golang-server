package core_http_utils

import (
	"net/http"
	"testing"
)

func TestRecorderDefaultsToOK(t *testing.T) {
	r := NewRecorder()
	if r.status != http.StatusOK {
		t.Fatalf("status = %d, want %d", r.status, http.StatusOK)
	}
	if !r.wroteHeader {
		r.Write([]byte("ok"))
	}
	if r.status != http.StatusOK {
		t.Fatalf("status after Write = %d", r.status)
	}
	if r.body.String() != "ok" {
		t.Fatalf("body = %q", r.body.String())
	}
}

func TestRecorderOnlyUsesFirstStatus(t *testing.T) {
	r := NewRecorder()
	r.WriteHeader(http.StatusCreated)
	r.WriteHeader(http.StatusBadRequest)
	if r.status != http.StatusCreated {
		t.Fatalf("status = %d, want %d", r.status, http.StatusCreated)
	}
}

func TestRecorderDoesNotStoreBodyForNoContent(t *testing.T) {
	r := NewRecorder()
	r.WriteHeader(http.StatusNoContent)
	n, err := r.Write([]byte("ignored"))
	if err != nil {
		t.Fatal(err)
	}
	if n != len("ignored") {
		t.Fatalf("Write returned %d", n)
	}
	if r.body.Len() != 0 {
		t.Fatalf("unexpected body: %q", r.body.String())
	}
}

func TestRecorderHeaderIsMutable(t *testing.T) {
	r := NewRecorder()
	r.Header().Set("Content-Type", "application/json")
	if got := r.header.Get("Content-Type"); got != "application/json" {
		t.Fatalf("Content-Type = %q", got)
	}
}
