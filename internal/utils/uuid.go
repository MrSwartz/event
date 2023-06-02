package utils

import "github.com/google/uuid"

func NewUUID() string {
	return uuid.New().String()
}

func CheckUUID(s string) bool {
	if _, err := uuid.Parse(s); err != nil {
		return false
	}
	return true
}
