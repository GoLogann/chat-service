package utils

import (
	"github.com/google/uuid"
)

func GenerateUUID() string {
	newUUID, err := uuid.NewV7()
	if err != nil {
		return "error"
	}
	return newUUID.String()
}
