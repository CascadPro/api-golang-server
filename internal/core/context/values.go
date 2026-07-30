package core_context

import (
	"context"
	"fmt"

	"github.com/CascadePro/api-golang-server/internal/core/domain"
	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

func UserIDFromContext(ctx context.Context) (uuid.UUID, error) {
	uidString, ok := ctx.Value(CtxKeyUserID).(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("`UserID` must be `string`: %w", core_errors.ErrInvalidArgument)
	}

	uid, err := uuid.Parse(uidString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("`UserID` must be `UUID`: %w", core_errors.ErrInvalidArgument)
	}

	return uid, nil
}

func UserRoleFromContext(ctx context.Context) (domain.UserRole, error) {
	r, ok := ctx.Value(CtxKeyUserRole).(domain.UserRole)
	if !ok {
		return domain.RoleRegular, fmt.Errorf("`Role` must be `UserRole`: %w", core_errors.ErrInvalidArgument)
	}

	return r, nil
}

func SessionIDFromContext(ctx context.Context) (string, error) {
	sid, ok := ctx.Value(CtxKeySessionID).(string)
	if !ok {
		return "", fmt.Errorf("`SessionID` must be `string`: %w", core_errors.ErrInvalidArgument)
	}
	if err := core_validation.ValidateID(sid, domain.SessionIDByteLength); err != nil {
		return "", fmt.Errorf("`SessionID` must be valid id: %w", err)
	}

	return sid, nil
}

func MimeTypeFromContext(ctx context.Context) (domain.FileMimeType, error) {
	m, ok := ctx.Value(CtxKeyMimeType).(domain.FileMimeType)
	if !ok {
		return domain.FileMimeTypeNil, fmt.Errorf("`MimeType` must be `FileMimeType`: %w", core_errors.ErrInvalidArgument)
	}

	return m, nil
}

func TagFromContext(ctx context.Context) (domain.FileTag, error) {
	t, ok := ctx.Value(CtxKeyTag).(domain.FileTag)
	if !ok {
		return domain.FileTagNil, fmt.Errorf("`Tag` must be `FileTag`: %w", core_errors.ErrInvalidArgument)
	}

	return t, nil
}
