package idempotency_service

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"sort"
	"strings"

	idempotency_error "github.com/CascadePro/api-golang-server/internal/features/idempotency/errors"
)

const maxRequestBodySize int64 = 12 << 20

type multipartPart struct {
	name        string
	filename    string
	contentType string
	contentHash [32]byte
}

func requestFingerprint(r *http.Request) ([]byte, error) {
	var body []byte

	if r.Body != nil {
		raw, err := io.ReadAll(io.LimitReader(r.Body, maxRequestBodySize+1))
		if err != nil {
			return nil, fmt.Errorf("read request body: %w", err)
		}

		if int64(len(raw)) > maxRequestBodySize {
			return nil, idempotency_error.ErrRequestBodyTooLarge
		}

		body = raw

		r.Body = io.NopCloser(bytes.NewReader(raw))
	}

	mediaType, params, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil && r.Header.Get("Content-Type") != "" {
		return nil, fmt.Errorf("parse content type: %w", err)
	}

	var canonicalBody []byte

	if strings.EqualFold(mediaType, "multipart/form-data") {
		canonicalBody, err = fingerprintMultipart(body, params["boundary"])
		if err != nil {
			return nil, err
		}
	} else {
		canonicalBody = body
	}

	hash := sha256.New()

	writePart(hash, r.Method)
	writePart(hash, r.URL.EscapedPath())
	writePart(hash, r.URL.RawQuery)
	writePart(hash, mediaType)
	writePart(hash, string(canonicalBody))

	return hash.Sum(nil), nil
}

func fingerprintMultipart(body []byte, boundary string) ([]byte, error) {
	if boundary == "" {
		return nil, fmt.Errorf("multipart boundary is missing")
	}

	reader := multipart.NewReader(bytes.NewReader(body), boundary)
	parts := make([]multipartPart, 0, 2)

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("read multipart part: %w", err)
		}

		content, readErr := io.ReadAll(io.LimitReader(part, maxRequestBodySize+1))
		closeErr := part.Close()

		if readErr != nil {
			return nil, fmt.Errorf("read multipart content: %w", readErr)
		}
		if closeErr != nil {
			return nil, fmt.Errorf("close multipart part: %w", closeErr)
		}

		if int64(len(content)) > maxRequestBodySize {
			return nil, idempotency_error.ErrRequestBodyTooLarge
		}

		parts = append(parts, multipartPart{
			name:        part.FormName(),
			filename:    part.FileName(),
			contentType: part.Header.Get("Content-Type"),
			contentHash: sha256.Sum256(content),
		})
	}

	sort.Slice(parts, func(i, j int) bool {
		a, b := parts[i], parts[j]

		if a.name != b.name {
			return a.name < b.name
		}

		if a.filename != b.filename {
			return a.filename < b.filename
		}

		return bytes.Compare(a.contentHash[:], b.contentHash[:]) < 0
	})

	hash := sha256.New()
	for _, part := range parts {
		writePart(hash, part.name)
		writePart(hash, part.filename)
		writePart(hash, part.contentType)

		_, _ = hash.Write(part.contentHash[:])
	}

	return hash.Sum(nil), nil
}

func writePart(w io.Writer, value string) {
	var length [8]byte

	binary.BigEndian.PutUint64(length[:], uint64(len(value)))

	_, _ = w.Write(length[:])
	_, _ = io.WriteString(w, value)
}
