package utils

import (
	"fmt"
	"regexp"
	"strings"
)

var SplitDIDRegexp, _ = regexp.Compile(`^did:([^:]+?)(:([^:]+?))?:([^:]+)$`)
var DidNamespaceRegexp, _ = regexp.Compile(`^[a-zA-Z0-9]*$`)

// Base58 only allowed (without OolI and 0)
var ValidBase58Regexp, _ = regexp.Compile(`^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$`)

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
		return "", "", "", fmt.Errorf("did must match the following regex exactly one time: %s", SplitDIDRegexp)
	}

	match := matches[0]
	return match[1], match[3], match[4], nil
}

// ValidateDID checks method and allowed namespaces only when the corresponding parameters are specified.
func ValidateDID(did string, method string, allowedNamespaces []string) error {
	sMethod, sNamespace, sUniqueId, err := TrySplitDID(did)
	if err != nil {
		return err
	}

	// check method
	// TODO: Is there any regexp fo method?
	if method != "" && method != sMethod {
		return fmt.Errorf("did method must be: %s", method)
	}

	// check namespaces
	if !DidNamespaceRegexp.MatchString(sNamespace) {
		return fmt.Errorf("did namespace must match the following regexp: %s", DidNamespaceRegexp)
	}

	if len(allowedNamespaces) > 0 && !Contains(allowedNamespaces, sNamespace) {
		return fmt.Errorf("did namespace must be one of: %s", strings.Join(allowedNamespaces[:], ", "))
	}

	// check unique-id
	err = ValidateUniqueId(sUniqueId)
	if err != nil {
		return err
	}

	return err
}

func ValidateUniqueId(uniqueId string) error {
	// Length should be 16 or 32 symbols
	if len(uniqueId) != 16 && len(uniqueId) != 32 {
		return fmt.Errorf("unique id length should be 16 or 32 symbols")
	}
	// Base58 check
	if !ValidBase58Regexp.MatchString(uniqueId) {
		return fmt.Errorf("unique id must be valid base58 string")
	}

	return nil
}

func IsValidDID(did string, method string, allowedNamespaces []string) bool {
	err := ValidateDID(did, method, allowedNamespaces)
	return err == nil
}
