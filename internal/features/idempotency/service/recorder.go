package idempotency_service

import (
	"bytes"
	"net/http"
)

type recorder struct {
	header      http.Header
	status      int
	wroteHeader bool
	body        bytes.Buffer
}

func newRecorder() *recorder {
	return &recorder{header: make(http.Header), status: http.StatusOK}
}

func (r *recorder) Header() http.Header { return r.header }

func (r *recorder) WriteHeader(status int) {
	if r.wroteHeader {
		return
	}
	if status < 100 || status > 999 {
		panic("invalid WriteHeader code")
	}
	r.status = status
	r.wroteHeader = true
}

func (r *recorder) Write(body []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}
	if r.status == http.StatusNoContent || r.status == http.StatusNotModified {
		return len(body), nil
	}
	return r.body.Write(body)
}
