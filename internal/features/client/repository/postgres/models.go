package client_postgres_repository

import (
	"encoding/json"
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type ClientModel struct {
	ID      uuid.UUID
	Version int64

	Company  string
	Contacts []string

	CreatedAt time.Time
}

type storedContact struct {
	Email *string `json:"email,omitempty"`
	Phone *string `json:"phone,omitempty"`
}

func domainFromModel(model ClientModel) domain.Client {
	contacts := make([]domain.ClientContact, 0, len(model.Contacts))
	for _, c := range model.Contacts {
		var sc storedContact
		if err := json.Unmarshal([]byte(c), &sc); err != nil {
			continue
		}

		var contact domain.ClientContact
		if sc.Email != nil {
			e := domain.Email(*sc.Email)
			contact.Email = &e
		}

		if sc.Phone != nil {
			p, err := domain.PhoneNumberFromE164(*sc.Phone)
			if err == nil {
				contact.Phone = &p
			}
		}

		if contact.Email != nil || contact.Phone != nil {
			contacts = append(contacts, contact)
		}
	}

	return domain.Client{
		ID:        model.ID,
		Version:   model.Version,
		Company:   model.Company,
		Contacts:  contacts,
		CreatedAt: model.CreatedAt,
	}
}

func domainsFromModels(models []ClientModel) []domain.Client {
	clients := make([]domain.Client, len(models))
	for i, model := range models {
		clients[i] = domainFromModel(model)
	}

	return clients
}

func modelFromDomain(client domain.Client) ClientModel {
	contacts := make([]string, 0, len(client.Contacts))
	for _, c := range client.Contacts {
		sc := storedContact{}
		if c.Email != nil {
			s := c.Email.String()
			sc.Email = &s
		}
		if c.Phone != nil {
			s := c.Phone.String()
			sc.Phone = &s
		}
		data, _ := json.Marshal(sc)
		contacts = append(contacts, string(data))
	}

	return ClientModel{
		ID:        client.ID,
		Version:   client.Version,
		Company:   client.Company,
		Contacts:  contacts,
		CreatedAt: client.CreatedAt,
	}
}
