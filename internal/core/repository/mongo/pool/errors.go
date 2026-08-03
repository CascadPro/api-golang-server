package core_mongo_pool

import (
	"errors"
	"fmt"

	core_errors "github.com/CascadePro/api-golang-server/internal/core/errors"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func MapErrors(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, mongo.ErrFileNotFound) {
		return fmt.Errorf("%v: %w", err, core_errors.ErrNotFound)
	}

	return fmt.Errorf("mongo operation failed: %w", err)
}
