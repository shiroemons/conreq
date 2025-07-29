// Package requestid provides utilities for generating and validating request IDs.
package requestid

import (
	"github.com/google/uuid"
)

// Generate creates a new UUID v4 request ID.
func Generate() string {
	return uuid.New().String()
}

// IsValid checks if the given string is a valid UUID.
func IsValid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
