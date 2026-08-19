package domain

import (
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
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

func (c *Client) Validate() error {
	if c.Company == "" {
		return fmt.Errorf("`Company` can't be NULL: %w", core_errors.ErrInvalidArgument)
	}
	if len(c.Contacts) == 0 {
		return fmt.Errorf("`Contacts` can't be EMPTY: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
