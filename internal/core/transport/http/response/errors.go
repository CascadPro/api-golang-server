package core_http_response

import "time"

type ErrorResponse struct {
	Code      string            `json:"code"             example:"invalid_credentials"`
	Message   string            `json:"message"          example:"short human-readable message"`
	Fields    map[string]string `json:"fields,omitempty"`
	Timestamp time.Time         `json:"timestamp"        example:"2006-01-02T15-04-05.000000"`
}
