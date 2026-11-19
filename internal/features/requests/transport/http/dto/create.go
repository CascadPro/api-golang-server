package request_http_dto

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type CreateRequestRequest struct {
	Title     string                         `json:"title"      validate:"required,min=1,max=255"        example:"Some title"`
	Origin    *[]CreateRequestRequestOrigin  `json:"origin"     validate:"omitempty,dive"`
	ClientID  *uuid.UUID                     `json:"client_id"  validate:"omitempty,uuid"                example:"00000000-0000-0000-0000-000000000000"`
	WorkTypes *[]string                      `json:"work_types" validate:"omitempty,dive,min=1,max=255"  example:""`
	Geography *[]string                      `json:"geography"  validate:"omitempty,dive,min=1,max=255"  example:""`
	Deadline  *time.Time                     `json:"deadline"   validate:"omitempty"                     example:"2006-01-02T15:04:05.000000Z"`
}

type CreateRequestRequestOrigin struct {
	Type  domain.RequestOriginType `json:"type"  validate:"required,oneof=email phone other" example:"email"`
	Email *string                  `json:"email,omitempty" validate:"omitempty,email"               example:"test@example.com"`
	Phone *string                  `json:"phone,omitempty" validate:"omitempty"                     example:"+79001234567"`
	Other *string                  `json:"other,omitempty" validate:"omitempty,min=1"               example:"Telegram: @user"`
}

func CreateRequestDomainFromDTO(dto CreateRequestRequest) domain.Request {
	request := domain.NewCreateRequest(dto.Title)

	if dto.Origin != nil {
		origin := make([]domain.RequestOrigin, 0, len(*dto.Origin))
		for _, j := range *dto.Origin {
			var o domain.RequestOrigin
			o.Type = j.Type

			switch j.Type {
			case domain.RequestOriginTypeEmail:
				if j.Email != nil {
					e, err := domain.ParseEmail(*j.Email)
					if err != nil {
						// Validation should catch this, but handle gracefully
						continue
					}
					o.Email = &e
				}
			case domain.RequestOriginTypePhone:
				if j.Phone != nil {
					p, err := domain.ParsePhoneNumber(*j.Phone)
					if err != nil {
						continue
					}
					o.Phone = &p
				}
			case domain.RequestOriginTypeOther:
				if j.Other != nil {
					o.Other = j.Other
				}
			}

			origin = append(origin, o)
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