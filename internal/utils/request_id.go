package utils

import (
	uuid "github.com/satori/go.uuid"
)

// Generate unique request id
func GenerateRequestID() string {
	return uuid.NewV4().String()
}
