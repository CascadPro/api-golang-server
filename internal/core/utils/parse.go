package core_utils

import (
	"context"
	"fmt"

	core_validation "github.com/CascadePro/api-golang-server/internal/core/validation"
	"github.com/google/uuid"
)

func ParseUUIDFromContext(ctx context.Context, name string) (uuid.UUID, error) {
	userID := ctx.Value(name)
	if userID == nil {
		return uuid.Nil, fmt.Errorf("`%s` can't be NULL", name)
	}

	userIDString, ok := userID.(string)
	if !ok {
		return uuid.Nil, fmt.Errorf("`%s` can't be parsed as string", name)
	}

	userUUID, err := uuid.Parse(userIDString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("`%s` can't be parsed as UUID: %w", name, err)
	}

	return userUUID, nil
}

func ParseIDStingFromContext(ctx context.Context, name string, length int) (string, error) {
	sessionIDAny := ctx.Value(name)
	if sessionIDAny == nil {
		return "", fmt.Errorf("`%s` can't be NULL", name)
	}

	sessionID, ok := sessionIDAny.(string)
	if !ok {
		return "", fmt.Errorf("`%s` can't be parsed as string", name)
	}

	if err := core_validation.ValidateID(sessionID, length); err != nil {
		return "", fmt.Errorf("`%s` isn't valid ID string: %w", name, err)
	}

	return sessionID, nil
}

func ParseFloat(s string) float64 {
	var f float64

	if _, err := fmt.Sscanf(s, "%f", &f); err != nil {
		return 0
	}

	return f
}
