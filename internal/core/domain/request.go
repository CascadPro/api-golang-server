package domain

import (
	"fmt"
	"slices"
	"time"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

type RequestStatus string
type RequestOriginType string

const (
	RequestStatusApproved = RequestStatus("approved")
	RequestStatusRejected = RequestStatus("rejected")
	RequestStatusWaiting  = RequestStatus("waiting")
	RequestStatusDefault  = RequestStatus("default")
	RequestStatusNil      = RequestStatus("")
)

const (
	RequestOriginTypeEmail = RequestOriginType("email")
	RequestOriginTypePhone = RequestOriginType("phone")
	RequestOriginTypeOther = RequestOriginType("other")
)

var (
	RequestStatuses = []RequestStatus{
		RequestStatusApproved, RequestStatusDefault, RequestStatusRejected, RequestStatusWaiting,
	}
	RequestOriginTypes = []RequestOriginType{
		RequestOriginTypeEmail, RequestOriginTypePhone, RequestOriginTypeOther,
	}
)

// ValidateRequestOriginType validates that the given origin type is one of the allowed values.
func ValidateRequestOriginType(t RequestOriginType) error {
	if slices.Contains(RequestOriginTypes, t) {
		return nil
	}
	return fmt.Errorf("invalid request origin type: %s", t)
}

type Request struct {
	ID      uuid.UUID
	Version int64
	Status  RequestStatus

	Title    string
	Origin   []RequestOrigin
	ClientID *uuid.UUID

	TechTaskDocID      *string
	ProjectDocID       *string
	SpecificationDocID *string

	WorkTypes []string
	Geography []string

	ContractDocID *string
	Deadline      *time.Time
	StatusBy      *uuid.UUID

	RequiredEmptyFields []string

	CreatedAt time.Time
	UpdatedAt time.Time
}

type RequestOrigin struct {
	Type  RequestOriginType
	Email *Email       `json:"email,omitempty"`
	Phone *PhoneNumber `json:"phone,omitempty"`
	Other *string      `json:"other,omitempty"`
}

// Validate ensures the origin has the correct value for its type.
func (o *RequestOrigin) Validate() error {
	switch o.Type {
	case RequestOriginTypeEmail:
		if o.Email == nil || o.Email.IsZero() {
			return fmt.Errorf("email origin must have a valid email: %w", core_errors.ErrInvalidArgument)
		}
		if o.Phone != nil || o.Other != nil {
			return fmt.Errorf("email origin cannot have phone or other: %w", core_errors.ErrInvalidArgument)
		}
	case RequestOriginTypePhone:
		if o.Phone == nil || o.Phone.IsZero() {
			return fmt.Errorf("phone origin must have a valid phone: %w", core_errors.ErrInvalidArgument)
		}
		if o.Email != nil || o.Other != nil {
			return fmt.Errorf("phone origin cannot have email or other: %w", core_errors.ErrInvalidArgument)
		}
	case RequestOriginTypeOther:
		if o.Other == nil || *o.Other == "" {
			return fmt.Errorf("other origin must have a non-empty value: %w", core_errors.ErrInvalidArgument)
		}
		if o.Email != nil || o.Phone != nil {
			return fmt.Errorf("other origin cannot have email or phone: %w", core_errors.ErrInvalidArgument)
		}
	default:
		return fmt.Errorf("unknown origin type: %w", core_errors.ErrInvalidArgument)
	}
	return nil
}

// GetValue returns the string representation of the origin value for backward compatibility.
func (o *RequestOrigin) GetValue() string {
	switch o.Type {
	case RequestOriginTypeEmail:
		if o.Email != nil {
			return o.Email.String()
		}
	case RequestOriginTypePhone:
		if o.Phone != nil {
			return o.Phone.String()
		}
	case RequestOriginTypeOther:
		if o.Other != nil {
			return *o.Other
		}
	}
	return ""
}

// NewRequestOriginEmail creates an email origin.
func NewRequestOriginEmail(email string) (RequestOrigin, error) {
	e, err := ParseEmail(email)
	if err != nil {
		return RequestOrigin{}, err
	}
	return RequestOrigin{
		Type:  RequestOriginTypeEmail,
		Email: &e,
	}, nil
}

// NewRequestOriginPhone creates a phone origin.
func NewRequestOriginPhone(phone string) (RequestOrigin, error) {
	p, err := ParsePhoneNumber(phone)
	if err != nil {
		return RequestOrigin{}, err
	}
	return RequestOrigin{
		Type:  RequestOriginTypePhone,
		Phone: &p,
	}, nil
}

// NewRequestOriginOther creates an other origin.
func NewRequestOriginOther(value string) RequestOrigin {
	return RequestOrigin{
		Type:  RequestOriginTypeOther,
		Other: &value,
	}
}

func NewCreateRequest(title string) Request {
	return Request{
		ID:      UninitializedUUID,
		Version: UninitializedVersion,
		Status:  RequestStatusDefault,
		Title:   title,
	}
}

func NewStatusPatchRequest(version int64, status RequestStatus, userID uuid.UUID) Request {
	return Request{
		Version:  version,
		Status:   status,
		StatusBy: &userID,
	}
}

func (r *Request) Validate() error {
	if err := core_validation.ValidateStringLength(&r.Title, "Title", 1, 255); err != nil {
		return err
	}

	if r.Status == RequestStatusNil {
		return fmt.Errorf("`Status` can't be NULL: %w`", core_errors.ErrInvalidArgument)
	}

	if r.Status == RequestStatusApproved {
		if len(r.Origin) == 0 {
			return fmt.Errorf("`Origin` can't be NULL if `Status` = 'approved': %w", core_errors.ErrInvalidArgument)
		}
		if len(r.WorkTypes) == 0 {
			return fmt.Errorf("`WorkTypes` can't be NULL if `Status` = 'approved': %w", core_errors.ErrInvalidArgument)
		}
		if len(r.Geography) == 0 {
			return fmt.Errorf("`Geography` can't be NULL if `Status` = 'approved': %w", core_errors.ErrInvalidArgument)
		}
		if r.ClientID == nil || *r.ClientID == uuid.Nil {
			return fmt.Errorf("`ClientID` can't be NULL if `Status` = 'approved': %w", core_errors.ErrInvalidArgument)
		}
		if err := r.ValidateDocs(); err != nil {
			return err
		}
	}

	// Validate each origin
	for i, origin := range r.Origin {
		if err := origin.Validate(); err != nil {
			return fmt.Errorf("origin[%d]: %w", i, err)
		}
	}

	return r.ValidateDate()
}

func (r *Request) ValidateDocs() error {
	if err := validateDocID(r.TechTaskDocID, "TechTaskDocID"); err != nil {
		return err
	}
	if err := validateDocID(r.ProjectDocID, "ProjectDocID"); err != nil {
		return err
	}
	if err := validateDocID(r.SpecificationDocID, "SpecificationDocID"); err != nil {
		return err
	}

	if r.ContractDocID != nil {
		if err := validateDocID(r.ContractDocID, "ContractDocID"); err != nil {
			return err
		}
	}

	return nil
}

func (r *Request) ValidateDate() error {
	if r.Deadline != nil && r.CreatedAt.After(*r.Deadline) {
		return fmt.Errorf("`Deadline` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}
	if r.CreatedAt.After(r.UpdatedAt) {
		return fmt.Errorf("`UpdatedAt` can't be before `CreatedAt`: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

func validateDocID(id *string, field string) error {
	if id == nil {
		return fmt.Errorf("`%s` can't be NULL if `Status` = 'approved': %w", field, core_errors.ErrInvalidArgument)
	}
	if err := core_validation.ValidateID(*id, FileIDByteLength); err != nil {
		return fmt.Errorf("`%s` is invalid file id: %w", field, err)
	}
	return nil
}

type RequestOriginPatch struct {
	Type  Nullable[RequestOriginType]
	Email Nullable[string]
	Phone Nullable[string]
	Other Nullable[string]
}

func (p *RequestOriginPatch) Validate() error {
	if p.Type.Set && p.Type.Value != nil {
		if err := ValidateRequestOriginType(*p.Type.Value); err != nil {
			return fmt.Errorf("invalid origin type: %w", err)
		}
	}

	valueFieldsSet := 0
	if p.Email.Set {
		valueFieldsSet++
	}
	if p.Phone.Set {
		valueFieldsSet++
	}
	if p.Other.Set {
		valueFieldsSet++
	}

	if valueFieldsSet > 1 {
		return fmt.Errorf("origin patch can only set one of email, phone, or other: %w", core_errors.ErrInvalidArgument)
	}

	if p.Email.Set && p.Email.Value != nil {
		if _, err := ParseEmail(*p.Email.Value); err != nil {
			return fmt.Errorf("invalid email: %w", err)
		}
	}
	if p.Phone.Set && p.Phone.Value != nil {
		if _, err := ParsePhoneNumber(*p.Phone.Value); err != nil {
			return fmt.Errorf("invalid phone: %w", err)
		}
	}
	if p.Other.Set && p.Other.Value != nil && *p.Other.Value == "" {
		return fmt.Errorf("other value cannot be empty: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}

type RequestPatch struct {
	Title              Nullable[string]
	Origin             Nullable[[]RequestOriginPatch]
	ClientID           Nullable[uuid.UUID]
	TechTaskDocID      Nullable[string]
	ProjectDocID       Nullable[string]
	SpecificationDocID Nullable[string]
	WorkTypes          Nullable[[]string]
	Geography          Nullable[[]string]
	ContractDocID      Nullable[string]
	Deadline           Nullable[time.Time]
}

func (p *RequestPatch) Validate() error {
	if p.Title.Set && p.Title.Value == nil {
		return fmt.Errorf("`Title` can't be patched to NULL: %w", core_errors.ErrInvalidArgument)
	}

	if p.Origin.Set {
		for i, origin := range *p.Origin.Value {
			if err := origin.Validate(); err != nil {
				return fmt.Errorf("origin[%d]: %w", i, err)
			}
		}
	}

	return nil
}

func (r *Request) ApplyPatch(patch RequestPatch) error {
	if err := patch.Validate(); err != nil {
		return fmt.Errorf("validate request patch: %w", err)
	}

	tmp := *r

	if patch.Title.Set {
		tmp.Title = *patch.Title.Value
	}

	if patch.Origin.Set {
		newOrigins := make([]RequestOrigin, 0, len(*patch.Origin.Value))
		for _, op := range *patch.Origin.Value {
			origin, err := op.ToRequestOrigin()
			if err != nil {
				return fmt.Errorf("convert origin patch: %w", err)
			}

			newOrigins = append(newOrigins, origin)
		}

		tmp.Origin = newOrigins
	}

	if patch.ClientID.Set {
		tmp.ClientID = patch.ClientID.Value
	}
	if patch.TechTaskDocID.Set {
		tmp.TechTaskDocID = patch.TechTaskDocID.Value
	}
	if patch.ProjectDocID.Set {
		tmp.ProjectDocID = patch.ProjectDocID.Value
	}
	if patch.SpecificationDocID.Set {
		tmp.SpecificationDocID = patch.SpecificationDocID.Value
	}
	if patch.WorkTypes.Set {
		if patch.WorkTypes.Value == nil {
			tmp.WorkTypes = nil
		} else {
			tmp.WorkTypes = *patch.WorkTypes.Value
		}
	}
	if patch.Geography.Set {
		if patch.Geography.Value == nil {
			tmp.Geography = nil
		} else {
			tmp.Geography = *patch.Geography.Value
		}
	}
	if patch.ContractDocID.Set {
		tmp.ContractDocID = patch.ContractDocID.Value
	}
	if patch.Deadline.Set {
		tmp.Deadline = patch.Deadline.Value
	}

	if err := tmp.Validate(); err != nil {
		return fmt.Errorf("validate request after applying patch: %w", err)
	}

	*r = tmp

	return nil
}

func (p *RequestOriginPatch) ToRequestOrigin() (RequestOrigin, error) {
	if err := p.Validate(); err != nil {
		return RequestOrigin{}, err
	}

	var origin RequestOrigin
	if p.Type.Set && p.Type.Value != nil {
		origin.Type = *p.Type.Value
	} else {
		if p.Email.Set && p.Email.Value != nil {
			origin.Type = RequestOriginTypeEmail
		} else if p.Phone.Set && p.Phone.Value != nil {
			origin.Type = RequestOriginTypePhone
		} else if p.Other.Set && p.Other.Value != nil {
			origin.Type = RequestOriginTypeOther
		} else {
			return RequestOrigin{}, fmt.Errorf("origin patch must have a type or value: %w", core_errors.ErrInvalidArgument)
		}
	}

	switch origin.Type {
	case RequestOriginTypeEmail:
		if p.Email.Set {
			if p.Email.Value == nil {
				return RequestOrigin{}, fmt.Errorf("email cannot be null: %w", core_errors.ErrInvalidArgument)
			}

			e, err := ParseEmail(*p.Email.Value)
			if err != nil {
				return RequestOrigin{}, err
			}

			origin.Email = &e
		}

	case RequestOriginTypePhone:
		if p.Phone.Set {
			if p.Phone.Value == nil {
				return RequestOrigin{}, fmt.Errorf("phone cannot be null: %w", core_errors.ErrInvalidArgument)
			}

			pn, err := ParsePhoneNumber(*p.Phone.Value)
			if err != nil {
				return RequestOrigin{}, err
			}

			origin.Phone = &pn
		}

	case RequestOriginTypeOther:
		if p.Other.Set {
			if p.Other.Value == nil {
				return RequestOrigin{}, fmt.Errorf("other cannot be null: %w", core_errors.ErrInvalidArgument)
			}
			if *p.Other.Value == "" {
				return RequestOrigin{}, fmt.Errorf("other value cannot be empty: %w", core_errors.ErrInvalidArgument)
			}

			origin.Other = p.Other.Value
		}
	}

	return origin, nil
}

type RequestQueryFilter struct {
	Status []RequestStatus
	Sort   SortType
}

func (f *RequestQueryFilter) Validate() error {
	if f.Sort != SortTypeNil {
		if err := core_validation.ValidateArray(SortTypes, f.Sort); err != nil {
			return fmt.Errorf("`SortType` isn't valid enum type: %w", err)
		}
	}

	if length := len(f.Status); length > 0 {
		if length > len(RequestStatuses) {
			return fmt.Errorf(
				"`Status` length can't be more than total existing statuses: %w",
				core_errors.ErrInvalidArgument,
			)
		}

		for _, s := range f.Status {
			if err := core_validation.ValidateArray(RequestStatuses, s); err != nil {
				return fmt.Errorf("`Status` isn't valid enum type: %w", err)
			}
		}
	}

	return nil
}

type RequestStats struct {
	Total    int64
	ByStatus map[RequestStatus]int64
}

func (r *RequestStats) Validate() error {
	if r.Total < 0 {
		return fmt.Errorf("`Total` can't be negative: %w", core_errors.ErrInvalidArgument)
	}

	for key, value := range r.ByStatus {
		if key == RequestStatusNil {
			return fmt.Errorf("key of `ByStatus` can't be NULL: %w", core_errors.ErrInvalidArgument)
		}

		if value < 0 {
			return fmt.Errorf("value of `ByStatus` can't be negative: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
