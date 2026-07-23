package core_http_response

import "time"

type ErrorResponse struct {
	Error     string    `json:"error" example:"full error text"`
	Message   string    `json:"message" example:"short human-readable message"`
	Timestamp time.Time `json:"timestamp" example:"2006-01-02T15-04-05.000000"`
}
