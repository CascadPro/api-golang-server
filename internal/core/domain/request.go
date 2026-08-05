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

