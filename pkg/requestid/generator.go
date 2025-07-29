package requestid

import (
	"github.com/google/uuid"
)

func Generate() string {
	return uuid.New().String()
}

func IsValid(id string) bool {
	_, err := uuid.Parse(id)
	return err == nil
}
