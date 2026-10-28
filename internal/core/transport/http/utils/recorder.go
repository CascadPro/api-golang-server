package core_http_utils

import (
	"bytes"
	"net/http"
)

type Recorder struct {
	header      http.Header
	wroteHeader bool

	status int
	body   bytes.Buffer
}

func NewRecorder() *Recorder {
	return &Recorder{
		header: make(http.Header),
		status: http.StatusOK,
	}
}

func (r *Recorder) Header() http.Header {
	return r.header
}

func (r *Recorder) GetBody() bytes.Buffer {
	return r.body
}

func (r *Recorder) GetStatus() int {
	return r.status
}

func (r *Recorder) WriteHeader(status int) {
	if r.wroteHeader {
		return
	}

	if status < 100 || status > 999 {
		panic("invalid WriteHeader code")
	}

	r.status = status
	r.wroteHeader = true
}

func (r *Recorder) Write(body []byte) (int, error) {
	if !r.wroteHeader {
		r.WriteHeader(http.StatusOK)
	}

	if r.status == http.StatusNoContent || r.status == http.StatusNotModified {
		return len(body), nil
	}

	return r.body.Write(body)
}
