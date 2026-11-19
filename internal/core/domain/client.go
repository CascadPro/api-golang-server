package domain

import (
	"encoding/json"
	"fmt"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"github.com/google/uuid"
)

type ClientContact struct {
	Email *Email       `json:"email,omitempty"`
	Phone *PhoneNumber `json:"phone,omitempty"`
}

func (c *ClientContact) Validate() error {
	hasEmail := c.Email != nil && !c.Email.IsZero()
	hasPhone := c.Phone != nil && !c.Phone.IsZero()

	if !hasEmail && !hasPhone {
		return fmt.Errorf("contact must have either email or phone: %w", core_errors.ErrInvalidArgument)
	}
	if hasEmail && hasPhone {
		return fmt.Errorf("contact cannot have both email and phone: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func NewClientContactEmail(email string) (ClientContact, error) {
	e, err := ParseEmail(email)
	if err != nil {
		return ClientContact{}, err
	}

	return ClientContact{Email: &e}, nil
}

func NewClientContactPhone(phone string) (ClientContact, error) {
	p, err := ParsePhoneNumber(phone)
	if err != nil {
		return ClientContact{}, err
	}
	return ClientContact{Phone: &p}, nil
}

type ClientContactPatch struct {
	Email Nullable[string]
	Phone Nullable[string]
}

func (p *ClientContactPatch) Validate() error {
	emailSet := p.Email.Set && p.Email.Value != nil
	emailNull := p.Email.Set && p.Email.Value == nil

	phoneSet := p.Phone.Set && p.Phone.Value != nil
	phoneNull := p.Phone.Set && p.Phone.Value == nil

	if emailSet && phoneSet {
		return fmt.Errorf("contact patch cannot set both email and phone: %w", core_errors.ErrInvalidArgument)
	}

	if emailSet {
		if _, err := ParseEmail(*p.Email.Value); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}
	}
	if phoneSet {
		if _, err := ParsePhoneNumber(*p.Phone.Value); err != nil {
			return fmt.Errorf("invalid phone: %w", err)
		}
	}

	_ = emailNull
	_ = phoneNull

	return nil
}

func (p *ClientContactPatch) ToContact() (ClientContact, error) {
	if err := p.Validate(); err != nil {
		return ClientContact{}, err
	}

	var contact ClientContact
	if p.Email.Set {
		if p.Email.Value == nil {
			contact.Email = nil
		} else {
			e, err := ParseEmail(*p.Email.Value)
			if err != nil {
				return ClientContact{}, err
			}
			contact.Email = &e
		}
	}

	if p.Phone.Set {
		if p.Phone.Value == nil {
			contact.Phone = nil
		} else {
			pn, err := ParsePhoneNumber(*p.Phone.Value)
			if err != nil {
				return ClientContact{}, err
			}
			contact.Phone = &pn
		}
	}

	return contact, nil
}

type Client struct {
	ID      uuid.UUID
	Version int64

	Company  string
	Contacts []ClientContact

	CreatedAt time.Time
}

func NewCreateClient(company string, contacts []ClientContact) Client {
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

	for i, contact := range c.Contacts {
		if err := contact.Validate(); err != nil {
			return fmt.Errorf("contact[%d]: %w", i, err)
		}
	}

	return nil
}

type jsonClientContact struct {
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

func (c ClientContact) MarshalJSON() ([]byte, error) {
	var emailStr, phoneStr *string
	if c.Email != nil {
		s := c.Email.String()
		emailStr = &s
	}

	if c.Phone != nil {
		s := c.Phone.String()
		phoneStr = &s
	}

	body := jsonClientContact{
		Email: emailStr,
		Phone: phoneStr,
	}

	return json.Marshal(&body)
}

// UnmarshalJSON implements custom JSON unmarshaling for ClientContact.
func (c *ClientContact) UnmarshalJSON(data []byte) error {
	var aux jsonClientContact
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	if aux.Email != nil {
		e, err := ParseEmail(*aux.Email)
		if err != nil {
			return err
		}
		c.Email = &e
	}

	if aux.Phone != nil {
		p, err := ParsePhoneNumber(*aux.Phone)
		if err != nil {
			return err
		}
		c.Phone = &p
	}

	return nil
}
