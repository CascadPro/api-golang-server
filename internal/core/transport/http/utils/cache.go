package core_http_utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"
)

type CacheControlOptionsType string

const (
	CacheControlTypePublic  CacheControlOptionsType = "public"
	CacheControlTypePrivate CacheControlOptionsType = "private"
)

type CacheControlOptions struct {
	Type                 CacheControlOptionsType
	MaxAge               time.Duration
	StaleWhileRevalidate time.Duration
	NoStore              bool
	NoCache              bool
	Immutable            bool
}

func SetCacheControl(rw http.ResponseWriter, options CacheControlOptions) {
	values := []string{}
	if options.Type != "" {
		values = append(values, string(options.Type))
	}
	if options.MaxAge > 0 {
		values = append(values, fmt.Sprintf("max-age=%d", int(options.MaxAge.Seconds())))
	}
	if options.StaleWhileRevalidate > 0 {
		values = append(values, fmt.Sprintf("stale-while-revalidate=%d", int(options.StaleWhileRevalidate.Seconds())))
	}
	if options.NoStore {
		values = append(values, "no-store")
	}
	if options.NoCache {
		values = append(values, "no-cache")
	}
	if options.Immutable {
		values = append(values, "immutable")
	}
	rw.Header().Set("Cache-Control", strings.Join(values, ", "))
}
