package utils

import (
	"fmt"
	"strings"
)

func ValidateDIDList(did_list []string, method string, allowedNamespaces []string) error {
	errors := []string{}
	for _, did := range did_list {
		error := ValidateDID(did, method, allowedNamespaces)
		if error != nil {
			errors = append(errors, error.Error())
		}
	}
	if len(errors) > 0 {
		return fmt.Errorf("Errors after the validation process of DID's list: %s.", strings.Join(errors, ", "))
	}
	return nil
}

func IsValidDIDList(did_list []string, method string, allowedNamespaces []string) bool {
	err := ValidateDIDList(did_list, method, allowedNamespaces)

	return nil == err
}
