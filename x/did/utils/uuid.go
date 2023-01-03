package utils

import (
	"errors"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

const StandardUUIDLength = 36

func ValidateUUID(u string) error {
	if len(u) != StandardUUIDLength {
		return errors.New("uuid must be of length " + strconv.Itoa(StandardUUIDLength) + " (in form of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)")
	}

	_, err := uuid.Parse(u)
	return err
}

func IsValidUUID(u string) bool {
	return ValidateUUID(u) == nil
}

// Normalization

func NormalizeUUID(uuid string) string {
	return strings.ToLower(uuid)
}
