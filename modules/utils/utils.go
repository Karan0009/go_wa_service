package utils

import (
	"fmt"
	"os"

	"github.com/google/uuid"
)

func GenerateUUID() uuid.UUID {
	return uuid.New()
}

func EnsureDirExists(path string) error {
	err := os.MkdirAll(path, os.ModePerm) // Creates directory with full permissions
	if err != nil {
		return fmt.Errorf("failed to create directory: %v", err)
	}
	return nil
}
