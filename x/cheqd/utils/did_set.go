package utils

import (
	"fmt"
	"strings"
)

func ValidateDIDSet(did_list []string, method string, allowedNamespaces []string) error {
	errors := []string{}
	if !IsDidSet(did_list) {
		return fmt.Errorf("There are not unic elements in the list: %s", strings.Join(did_list, ", "))
	}
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

func IsDidSet(did_list []string) bool {
	set :=map[string]int{}
	for i, did := range did_list {
		set[did] = i
	}

	return len(did_list) == len(set)
}

func IsValidDIDSet(did_list []string, method string, allowedNamespaces []string) bool {
	err := ValidateDIDSet(did_list, method, allowedNamespaces)

	return nil == err
}
