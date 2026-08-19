package client_postgres_repository

import (
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

func domainFromModel(model ClientModel) domain.Client {
	return domain.Client{
		ID:        model.ID,
		Version:   model.Version,
		Company:   model.Company,
		Contacts:  model.Contacts,
		CreatedAt: model.CreatedAt,
	}
}
