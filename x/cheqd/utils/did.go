package utils

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	SplitDIDRegexp, _     = regexp.Compile(`^did:([^:]+?)(:([^:]+?))?:([^:]+)$`)
	DidNamespaceRegexp, _ = regexp.Compile(`^[a-zA-Z0-9]*$`)
)

// TrySplitDID Validates generic format of DID. It doesn't validate method, name and id content.
// Call ValidateDID for further validation.
func TrySplitDID(did string) (method string, namespace string, id string, err error) {
	// Example: did:cheqd:testnet:base58str1ng1111
	// match [0] - the whole string
	// match [1] - cheqd                - method
	// match [2] - :testnet
	// match [3] - testnet              - namespace
	// match [4] - base58str1ng1111     - id
	matches := SplitDIDRegexp.FindAllStringSubmatch(did, -1)
	if len(matches) != 1 {
		return "", "", "", errors.New("unable to split did into method, namespace and id")
	}

	match := matches[0]
	return match[1], match[3], match[4], nil
}

func MustSplitDID(did string) (method string, namespace string, id string) {
	method, namespace, id, err := TrySplitDID(did)
	if err != nil {
		panic(err.Error())
	}
	return
}

func JoinDID(method, namespace, id string) string {
	res := "did:" + method

	if namespace != "" {
		res = res + ":" + namespace
	}

	return res + ":" + id
}

// ValidateDID checks method and allowed namespaces only when the corresponding parameters are specified.
func ValidateDID(did string, method string, allowedNamespaces []string) error {
	sMethod, sNamespace, sUniqueId, err := TrySplitDID(did)
	if err != nil {
		return err
	}

	// check method
	if method != "" && method != sMethod {
		return fmt.Errorf("did method must be: %s", method)
	}

	// check namespaces
	if !DidNamespaceRegexp.MatchString(sNamespace) {
		return errors.New("invalid did namespace")
	}

	if len(allowedNamespaces) > 0 && !Contains(allowedNamespaces, sNamespace) {
		return fmt.Errorf("did namespace must be one of: %s", strings.Join(allowedNamespaces[:], ", "))
	}

	// check unique-id
	err = ValidateID(sUniqueId)
	if err != nil {
		return err
	}

	return err
}

func IsValidDID(did string, method string, allowedNamespaces []string) bool {
	err := ValidateDID(did, method, allowedNamespaces)
	return err == nil
}

func ValidateID(id string) error {
	isValidId := len(id) >= 16 && len(id) <= 32 && IsValidBase58(id) ||
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
