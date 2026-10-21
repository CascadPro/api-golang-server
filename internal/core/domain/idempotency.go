package domain

import "net/http"

type RecordStatus string

type IdempotencyRecord struct {
	Scope     string
	Operation string

	KeyHash     []byte
	RequestHash []byte

	Status          RecordStatus
	ResponseStatus  int
	ResponseBody    []byte
	ResponseHeaders http.Header
}

const (
	IdempotencyRecordStatusProcessing = RecordStatus("processing")
	IdempotencyRecordStatusCompleted  = RecordStatus("completed")
)
