package domain

import (
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ID      uuid.UUID
	Version int64

	Company  string
	Contacts []string

	CreatedAt time.Time
}

func NewCreateClient(company string, contacts []string) Client {
	return Client{
		ID:       UninitializedUUID,
		Version:  UninitializedVersion,
		Company:  company,
		Contacts: contacts,
	}
}
