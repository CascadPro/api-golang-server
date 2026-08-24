package requests_mongo_repository

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type RequestModel struct {
	ID      string               `bson:"_id"`
	Version int64                `bson:"v"`
	Status  domain.RequestStatus `bson:"status"`

	Title    string               `bson:"title"`
	Origin   []RequestModelOrigin `bson:"origin,omitempty"`
	ClientID *string              `bson:"client,omitempty"`

	Docs RequestModelDocs `bson:"docs"`

	WorkTypes []string `bson:"work_types,omitempty"`
	Geography []string `bson:"geo_desc,omitempty"`

	ContractDocID *string    `bson:"contract,omitempty"`
	Deadline      *time.Time `bson:"deadline,omitempty"`
	StatusBy      *string    `bson:"status_by,omitempty"`

	RequiredEmptyFields []string `bson:"required_empty_fields"`

	CreatedAt time.Time `bson:"created_at"`
	UpdatedAt time.Time `bson:"updated_at"`
}

type RequestModelOrigin struct {
	Type  domain.RequestOriginType `bson:"type"`
	Value string                   `bson:"value"`
}

type RequestModelDocs struct {
	TechTaskDocID      *string `bson:"tech_task,omitempty"`
	ProjectDocID       *string `bson:"project,omitempty"`
	SpecificationDocID *string `bson:"specification,omitempty"`
}

func (m *RequestModel) CalcRequiredEmptyFields() {
	fields := make([]string, 0, 7)

	if len(m.Origin) == 0 {
		fields = append(fields, "origin")
	}
	if len(m.WorkTypes) == 0 {
		fields = append(fields, "work_types")
	}
	if len(m.Geography) == 0 {
		fields = append(fields, "geo_desc")
	}
	if m.ClientID == nil || *m.ClientID == "" {
		fields = append(fields, "client")
	}
	if m.Docs.ProjectDocID == nil {
		fields = append(fields, "docs.project")
	}
	if m.Docs.TechTaskDocID == nil {
		fields = append(fields, "docs.tech_task")
	}
	if m.Docs.SpecificationDocID == nil {
		fields = append(fields, "docs.specification")
	}

	m.RequiredEmptyFields = fields
}

func domainToModel(request domain.Request) RequestModel {
	model := RequestModel{
		Version:   request.Version,
		Status:    request.Status,
		Title:     request.Title,
		Origin:    domainOriginsToModel(request.Origin),
		WorkTypes: request.WorkTypes,
		Geography: request.Geography,
		Docs: RequestModelDocs{
			TechTaskDocID:      request.TechTaskDocID,
			ProjectDocID:       request.ProjectDocID,
			SpecificationDocID: request.SpecificationDocID,
		},
		ContractDocID: request.ContractDocID,
		Deadline:      request.Deadline,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if request.ClientID != nil && *request.ClientID != uuid.Nil {
		clientStr := request.ClientID.String()
		model.ClientID = &clientStr
	}
	if request.StatusBy != nil && *request.StatusBy != uuid.Nil {
		statusByStr := request.StatusBy.String()
		model.StatusBy = &statusByStr
	}

	return model
}

func domainOriginsToModel(origin []domain.RequestOrigin) []RequestModelOrigin {
	if len(origin) == 0 {
		return nil
	}

	slice := make([]RequestModelOrigin, 0, len(origin))
	for _, j := range origin {
		slice = append(slice, RequestModelOrigin(j))
	}
	return slice
}

func modelToDomain(model RequestModel) domain.Request {
	origin := make([]domain.RequestOrigin, 0, len(model.Origin))
	for _, j := range model.Origin {
		origin = append(origin, domain.RequestOrigin(j))
	}

	requestDomain := domain.Request{
		ID:                  uuid.MustParse(model.ID),
		Version:             model.Version,
		Status:              model.Status,
		Title:               model.Title,
		Origin:              origin,
		WorkTypes:           model.WorkTypes,
		Geography:           model.Geography,
		ContractDocID:       model.ContractDocID,
		RequiredEmptyFields: model.RequiredEmptyFields,
		Deadline:            model.Deadline,
		CreatedAt:           model.CreatedAt,
		UpdatedAt:           model.UpdatedAt,
	}

	if model.ClientID != nil && *model.ClientID != "" {
		id := uuid.MustParse(*model.ClientID)
		requestDomain.ClientID = &id
	}
	if model.StatusBy != nil && *model.StatusBy != "" {
		id := uuid.MustParse(*model.StatusBy)
		requestDomain.StatusBy = &id
	}

	if model.Docs.ProjectDocID != nil && *model.Docs.ProjectDocID != "" {
		requestDomain.ProjectDocID = model.Docs.ProjectDocID
	}
	if model.Docs.SpecificationDocID != nil && *model.Docs.SpecificationDocID != "" {
		requestDomain.SpecificationDocID = model.Docs.SpecificationDocID
	}
	if model.Docs.TechTaskDocID != nil && *model.Docs.TechTaskDocID != "" {
		requestDomain.TechTaskDocID = model.Docs.TechTaskDocID
	}

	return requestDomain
}

func modelsToDomains(models []RequestModel) []domain.Request {
	requests := make([]domain.Request, 0, len(models))
	for _, model := range models {
		requests = append(requests, modelToDomain(model))
	}
	return requests
}
