package utils

import (
	"errors"
	"strings"
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

func NormalizeIdentifier(did string) string {
	_, _, sUniqueId, err := TrySplitDID(did)
	if err != nil {
		sUniqueId = did
	}
	if IsValidUUID(sUniqueId) {
		return strings.ToLower(did)
	}
	return did
}

func NormalizeIdForFragmentUrl(didUrl string) string {
	id, _, _, fragmentId, err := TrySplitDIDUrl(didUrl)
	if err != nil {
		return didUrl
	}
	if fragmentId == "" {
		return NormalizeIdentifier(id)
	}
	return NormalizeIdentifier(id) + "#" + fragmentId
}

func NormalizeIdentifiersList(keys []string) []string {
	if keys == nil {
		return nil
	}
	newKeys := []string{}
	for _, id := range keys {
		newKeys = append(newKeys, NormalizeIdForFragmentUrl(id))
	}
	return newKeys
}
