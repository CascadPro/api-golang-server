package idempotency_service

import (
	"bytes"
	"net/http"
)

type recorder struct {
	header http.Header
	status int
	body   bytes.Buffer
}

func newRecorder() *recorder {
	return &recorder{
		header: make(http.Header),
		status: http.StatusOK,
	}
}

func (r *recorder) Header() http.Header {
	return r.header
}

func (r *recorder) WriteHeader(status int) {
	if r.status != http.StatusOK {
		return
	}

	r.status = status
}

func (r *recorder) Write(body []byte) (int, error) {
	if r.body.Len() == 0 && r.status == http.StatusOK {
		r.status = http.StatusOK
	}

	return r.body.Write(body)
}
