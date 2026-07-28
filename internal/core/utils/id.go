package core_utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateID(length int) (string, error) {
	bytes := make([]byte, length)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("read bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
