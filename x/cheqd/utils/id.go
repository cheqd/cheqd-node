package utils

import (
	"errors"

	"github.com/mr-tron/base58"
)

const (
	IndyIdLength = 16
)

func ValidateID(id string) error {
	isValidId := IsValidIndyId(id) || IsValidUUID(id)

	if !isValidId {
		return errors.New("unique id should be one of: 16 bytes of decoded base58 string or UUID")
	}

	return nil
}

func IsValidID(id string) bool {
	err := ValidateID(id)
	return err == nil
}

func IsValidIndyId(data string) bool {
	bytes, err := base58.Decode(data)
	if err != nil {
		return false
	}
	return len(bytes) == IndyIdLength
}

// Normalization

func NormalizeId(id string) string {
	if IsValidUUID(id) {
		return NormalizeUUID(id)
	}
	return id
}

func NormalizeIdList(keys []string) []string {
	if keys == nil {
		return nil
	}
	newKeys := []string{}
	for _, id := range keys {
		newKeys = append(newKeys, NormalizeId(id))
	}
	return newKeys
}
