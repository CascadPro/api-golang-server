package request_http_dto

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type CreateRequestRequest struct {
	Title     string                        `json:"title"      validate:"required,min=1,max=255"        example:"Some title"`
	Origin    *[]CreateRequestRequestOrigin `json:"origin"     validate:"omitempty"                     example:""`
	ClientID  *uuid.UUID                    `json:"client_id"  validate:"omitempty,uuid"                example:"00000000-000000-000000-000000000000"`
	WorkTypes *[]string                     `json:"work_types" validate:"omitempty,dive,min=1,max=255"  example:""`
	Geography *[]string                     `json:"geography"  validate:"omitempty,dive,min=1,max=255"  example:""`
	Deadline  *time.Time                    `json:"deadline"   validate:"omitempty,datetime=2006-01-02" example:"2006-01-02T15-04-05.000000"`
}

type CreateRequestRequestOrigin struct {
	Type  domain.RequestOriginType `json:"type"  validate:"required,oneof=email phone other" example:"email"`
	Value string                   `json:"value" validate:"required,alphanumspace"           example:"test@example.com"`
}

func CreateRequestDomainFromDTO(dto CreateRequestRequest) domain.Request {
	request := domain.NewCreateRequest(dto.Title)

	if dto.Origin != nil {
		origin := make([]domain.RequestOrigin, len(*dto.Origin))
		for i, j := range *dto.Origin {
			origin[i] = domain.RequestOrigin(j)
		}

		request.Origin = origin
	}
	if dto.ClientID != nil {
		request.ClientID = dto.ClientID
	}
	if dto.WorkTypes != nil {
		request.WorkTypes = *dto.WorkTypes
	}
	if dto.Geography != nil {
		request.Geography = *dto.Geography
	}
	if dto.Deadline != nil {
		request.Deadline = dto.Deadline
	}

	return request
}
