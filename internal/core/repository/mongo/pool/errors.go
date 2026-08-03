package core_mongo_pool

import (
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/mongo"
)

var (
	ErrNotFound = errors.New("mongo not found")
	ErrUnknown  = errors.New("mongo unknown error")
)

func MapErrors(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, mongo.ErrNoDocuments) || errors.Is(err, mongo.ErrFileNotFound) {
		return fmt.Errorf("%v: %w", err, ErrNotFound)
	}

	return fmt.Errorf("mongo operation failed: %v: %w", err, ErrUnknown)
}
