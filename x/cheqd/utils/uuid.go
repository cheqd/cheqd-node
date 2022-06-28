package utils

import (
	"errors"
	"strconv"

	"github.com/google/uuid"
)

const StandardUuidLength = 36

func ValidateUUID(u string) error {
	if len(u) != StandardUuidLength {
		return errors.New("uuid must be of length " + strconv.Itoa(StandardUuidLength) + " (in form of xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx)")
	}

	_, err := uuid.Parse(u)
	return err
}

func IsValidUUID(u string) bool {
	return ValidateUUID(u) == nil
}
