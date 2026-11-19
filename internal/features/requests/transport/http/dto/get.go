package request_http_dto

import (
	"time"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	"github.com/google/uuid"
)

type GetRequestResponse struct {
	GetRequestsResponseRequest
	WorkTypes           []string `json:"work_types,omitempty"`
	Geography           []string `json:"geography_description,omitempty"`
	RequiredEmptyFields []string `json:"fields_remaining"`

	Client   *GetRequestResponseClient `json:"client"`
	Docs     *GetRequestResponseDocs   `json:"docs,omitempty"`
	StatusBy *GetRequestResponseUser   `json:"status_by,omitempty"`
}

type GetRequestResponseUser struct {
	ID           uuid.UUID       `json:"id"                       example:"00000000-0000-0000-0000-000000000000"`
	Role         domain.UserRole `json:"role"                     example:"regular"`
	Name         string          `json:"name"                     example:"John"`
	Surname      string          `json:"surname"                  example:"Doe"`
	AvatarFileID *string         `json:"avatar_file_id,omitempty" example:"831f21c1798a7972fa9cda12dac0"`
}

type GetRequestResponseClient struct {
	ID       uuid.UUID         `json:"id"      example:"00000000-0000-0000-0000-000000000000"`
	Company  string            `json:"company" example:"Cascade Pro"`
	Contacts []ClientContactResponse `json:"contacts"`
}

type ClientContactResponse struct {
	Email *string `json:"email,omitempty" example:"test@example.com"`
	Phone *string `json:"phone,omitempty" example:"+79001234567"`
}

type GetRequestResponseDocs struct {
	Project       *GetRequestResponseDocsFile `json:"project"`
	TechTask      *GetRequestResponseDocsFile `json:"tech_task"`
	Specification *GetRequestResponseDocsFile `json:"specification"`
	Contract      *GetRequestResponseDocsFile `json:"contract"`
}

type GetRequestResponseDocsFile struct {
	ID        string         `json:"id"         example:"831f21c1798a7972fa9cda12dac0"`
	Tag       domain.FileTag `json:"tag"        example:"avatars"`
	Filename  string         `json:"filename"   example:"avatar.jpg"`
	Size      int64          `json:"size"       example:"1024"`
	CreatedAt time.Time      `json:"created_at" example:"2006-01-02T15:04:05.000000Z"`
}

func RequestResponseFromDomain(request domain.Request, user domain.User, client domain.Client, files *map[int]domain.File) GetRequestResponse {
	var response GetRequestResponse

	response.GetRequestsResponseRequest = RequestsResponseFromDomain(request)

	response.WorkTypes = request.WorkTypes
	response.Geography = request.Geography
	response.RequiredEmptyFields = request.RequiredEmptyFields

	if request.StatusBy != nil && *request.StatusBy != uuid.Nil {
		response.StatusBy = &GetRequestResponseUser{
			ID:           user.ID,
			Role:         user.Role,
			Name:         user.Name,
			Surname:      user.Surname,
			AvatarFileID: user.AvatarFileID,
		}
	}

	if request.ClientID != nil && *request.ClientID != uuid.Nil {
		contacts := make([]ClientContactResponse, len(client.Contacts))
		for i, c := range client.Contacts {
			var cc ClientContactResponse
			if c.Email != nil {
				s := c.Email.String()
				cc.Email = &s
			}
			if c.Phone != nil {
				s := c.Phone.String()
				cc.Phone = &s
			}
			contacts[i] = cc
		}
		response.Client = &GetRequestResponseClient{
			ID:       client.ID,
			Company:  client.Company,
			Contacts: contacts,
		}
	}

	var docs GetRequestResponseDocs
	if request.ProjectDocID != nil {
		docs.Project = docsFileResponseFromDomain((*files)[0])
	}
	if request.TechTaskDocID != nil {
		docs.TechTask = docsFileResponseFromDomain((*files)[1])
	}
	if request.SpecificationDocID != nil {
		docs.Specification = docsFileResponseFromDomain((*files)[2])
	}
	if request.ContractDocID != nil {
		docs.Contract = docsFileResponseFromDomain((*files)[3])
	}

	response.Docs = &docs

	return response
}

func docsFileResponseFromDomain(file domain.File) *GetRequestResponseDocsFile {
	return &GetRequestResponseDocsFile{
		ID:        file.ID,
		Tag:       file.Tag,
		Filename:  file.Filename,
		Size:      int64(file.Size),
		CreatedAt: file.CreatedAt,
	}
}