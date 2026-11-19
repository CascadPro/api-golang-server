package request_http_dto

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type GetRequestsRequest struct {
	Limit    int                    `json:"limit"    validate:"omitempty,number,gte=1,lte=100" example:"10"`
	Page     int                    `json:"page"     validate:"omitempty,number,gte=1,lte=100" example:"1"`
	Sort     domain.SortType        `json:"sort"     validate:"omitempty,alpha,sort_type"      example:"newest"`
	Statuses []domain.RequestStatus `json:"statuses" validate:"omitempty,max=4,unique"`
}

type GetRequestsResponse struct {
	Requests    []GetRequestsResponseRequest `json:"requests"`
	TotalLength int64                        `json:"length" example:"1"`
}

type GetRequestsResponseRequest struct {
	ID       uuid.UUID            `json:"id"                  example:"00000000-0000-0000-0000-000000000000"`
	Status   domain.RequestStatus `json:"status"              example:"default"`
	StatusBy *uuid.UUID           `json:"status_by,omitempty" example:"00000000-0000-0000-0000-000000000000"`

	Title  string                      `json:"title"            example:"New Request"`
	Origin []GetRequestsResponseOrigin `json:"origin,omitempty"`

	ClientID *uuid.UUID `json:"client,omitempty"           example:"00000000-0000-0000-0000-000000000000"`

	RequiredEmptyFields int `json:"fields_remaining" example:"3"`

	Deadline  *time.Time `json:"deadline,omitempty"  example:"2006-01-02T15:04:05.000000Z"`
	CreatedAt time.Time  `json:"created_at"          example:"2006-01-02T15:04:05.000000Z"`
	UpdatedAt time.Time  `json:"updated_at"          example:"2006-01-02T15:04:05.000000Z"`
}

type GetRequestsResponseOrigin struct {
	Type  domain.RequestOriginType `json:"type"  example:"email"`
	Email *string                  `json:"email,omitempty" example:"test@example.com"`
	Phone *string                  `json:"phone,omitempty" example:"+79001234567"`
	Other *string                  `json:"other,omitempty" example:"Telegram: @user"`
}

func RequestsResponseFromDomain(request domain.Request) GetRequestsResponseRequest {
	response := GetRequestsResponseRequest{
		ID:                  request.ID,
		Status:              request.Status,
		StatusBy:            request.StatusBy,
		Title:               request.Title,
		ClientID:            request.ClientID,
		RequiredEmptyFields: len(request.RequiredEmptyFields),
		Deadline:            request.Deadline,
		CreatedAt:           request.CreatedAt,
		UpdatedAt:           request.UpdatedAt,
	}

	if len(request.Origin) > 0 {
		origins := make([]GetRequestsResponseOrigin, len(request.Origin))
		for i, origin := range request.Origin {
			origins[i] = GetRequestsResponseOriginFromDomain(origin)
		}
		response.Origin = origins
	}

	return response
}

func GetRequestsResponseOriginFromDomain(origin domain.RequestOrigin) GetRequestsResponseOrigin {
	resp := GetRequestsResponseOrigin{Type: origin.Type}
	switch origin.Type {
	case domain.RequestOriginTypeEmail:
		if origin.Email != nil {
			s := origin.Email.String()
			resp.Email = &s
		}
	case domain.RequestOriginTypePhone:
		if origin.Phone != nil {
			s := origin.Phone.String()
			resp.Phone = &s
		}
	case domain.RequestOriginTypeOther:
		if origin.Other != nil {
			resp.Other = origin.Other
		}
	}
	return resp
}