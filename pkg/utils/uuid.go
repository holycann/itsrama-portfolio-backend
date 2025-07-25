package utils

import "github.com/google/uuid"

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.NewString()
}

// GenerateUUIDIfEmpty generates a new UUID if the provided ID is empty
func GenerateUUIDIfEmpty(id string) string {
	if id == "" {
		return uuid.NewString()
	}
	return id
}
