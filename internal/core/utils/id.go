package core_utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
)

func GenerateID(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)

	_, err := rand.Read(bytes)
	if err != nil {
		return "", fmt.Errorf("read bytes: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
