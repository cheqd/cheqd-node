package utils

import (
	"errors"
)

func ValidateID(id string) error {
	isValidId := len(id) == 16 && IsValidBase58(id) ||
		len(id) == 32 && IsValidBase58(id) ||
		IsValidUUID(id)

	if !isValidId {
		return errors.New("unique id should be one of: 16 symbols base58 string, 32 symbols base58 string, or UUID")
	}

	return nil
}

func IsValidID(id string) bool {
	err := ValidateID(id)
	return err == nil
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
