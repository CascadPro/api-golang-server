package domain

import (
	"fmt"
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

type Request struct {
	ID      uuid.UUID
	Version int
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

	CreatedAt time.Time
	UpdatedAt time.Time
}

type RequestOrigin struct {
	Type  RequestOriginType
	Value string
}

func NewCreateRequest(title string) Request {
	return Request{
		ID:      UninitializedUUID,
		Version: UninitializedVersion,
		Status:  RequestStatusDefault,
		Title:   title,
	}
}

func NewStatusPatchRequest(version int, status RequestStatus, userID uuid.UUID) Request {
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

type RequestPatch struct {
	Title              Nullable[string]
	Origin             Nullable[[]RequestOrigin]
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
		tmp.Origin = *patch.Origin.Value
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
		tmp.WorkTypes = *patch.WorkTypes.Value
	}
	if patch.Geography.Set {
		tmp.Geography = *patch.Geography.Value
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

	if length := len(f.Status); length > 0 && length < len(RequestStatuses) {
		return fmt.Errorf("`Status` length can't be more than total existing statuses: %w", core_errors.ErrInvalidArgument)
	}

	return nil
}
