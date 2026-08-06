package core_http

import "net/http"

type Method string

const (
	MethodGet     = Method(http.MethodGet)
	MethodQuery   = Method("QUERY")
	MethodPost    = Method(http.MethodPost)
	MethodPatch   = Method(http.MethodPatch)
	MethodDelete  = Method(http.MethodDelete)
	MethodOptions = Method(http.MethodOptions)
)

var (
	Methods        = []Method{MethodGet, MethodQuery, MethodPost, MethodPatch, MethodDelete, MethodOptions}
	AllowedMethods = []Method{MethodGet, MethodQuery, MethodPost, MethodPatch, MethodDelete, MethodOptions}
)
