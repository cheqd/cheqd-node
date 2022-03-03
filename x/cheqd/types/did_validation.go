package types

import (
	"errors"
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	"regexp"
)

var SplitDIDRegexp, _        = regexp.Compile(`^did:([^:]+?)(:([^:]+?))?:([^:]+)$`)
var DidNamespaceRegexp, _    = regexp.Compile(`^[a-zA-Z0-9]$`)
// Base58 only allowed (without OolI and 0)
var UniqueIDRegexp, _        = regexp.Compile(`^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$`)
// That for groups:
// Example: did:cheqd:testnet:fafdsffq11213343/path-to-s/ome-external-resource?query#key1???
// 1 - [^/?#]* - all the symbols except / and ? and # . This is the DID part                      (did:cheqd:testnet:fafdsffq11213343)
// 2 - [^?#]*  - all the symbols except ? and #. it means te section started from /, path-abempty (/path-to-s/ome-external-resource)
// 3 - \?([^#]*) - group for `query` part but with ? symbol 									  (?query)
// 4 - [^#]*     - group inside query string, match only exact query                              (query)
// 5 - #([^#]+[\$]?) - group for fragment, starts with #, includes #                              (#key1???)
// 6 - [^#]+[\$]?    - fragment only															  (key1???)
// Number of queries is not limited.
var SplitDIDURL, _           = regexp.Compile(`([^/?#]*)?([^?#]*)(\?([^#]*))?(#([^#]+$))?$`)
var DIDPathAbemptyRegexp, _  = regexp.Compile(`^[/a-zA-Z0-9\-\.\_\~\%\!\$\&\'\(\)\*\+\,\;\=\:\@]+$`)
var DIDQueryRegexp, _        = regexp.Compile(`^[/a-zA-Z0-9\-\.\_\~\%\!\$\&\'\(\)\*\+\,\;\=\:\@\/\?]*$`)
var DIDFragmentRegexp, _     = regexp.Compile(`^[/a-zA-Z0-9\-\.\_\~\%\!\$\&\'\(\)\*\+\,\;\=\:\@\/\?]*$`)


//// DID-related

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
		return "", "", "", errors.New("there should be exactly one match")
	}

	match := matches[0]
	return match[1], match[3], match[4], nil
}

// ValidateDID checks method and allowed namespaces only when the corresponding parameters are specified.
// TODO: Handle empty parameters
func ValidateDID(did string, method string, allowedNamespaces []string) error {
	method, namespace, unique_id, err := TrySplitDID(did)
	if err != nil {
		return err
	}

	// check method
	if method != method {
		return ErrStaticDIDBadMethod.Wrap(method)
	}
	// check namespaces
	if !DidNamespaceRegexp.MatchString(namespace) || !utils.Contains(allowedNamespaces, namespace) {
		return ErrStaticDIDNamespaceNotAllowed.Wrap(namespace)
	}
	// check unique-id
	err := ValidateUniqueId(unique_id)

	return err
}

func ValidateUniqueId(unique_id string) error {
	// Length should be 16 or 32 symbols
	if len(unique_id) != 16 && len(unique_id) != 32 {
		return ErrStaticDIDBadUniqueIDLen.Wrap(unique_id)
	}
	// Base58 check
	if !UniqueIDRegexp.MatchString(unique_id) {
		return ErrStaticDIDNotBase58ID.Wrap(unique_id)
	}

	return nil
}

func IsValidDID(did string, method string, allowedNamespaces []string) bool {
	err := ValidateDID(did, method, allowedNamespaces)
	return err == nil
}


//// DID URL-related

// TrySplitDIDUrl Validates generic format of DIDUrl. It doesn't validate path, query and fragment content.
// Call ValidateDIDUrl for further validation.
func TrySplitDIDUrl(didUrl string) (did string, path string , query string, fragment string) {
	match := SplitDIDURL.FindStringSubmatch(didUrl)
	return match[1], match[2], match[4], match[6]
}

// ValidateDID checks method and allowed namespaces only when the corresponding parameters are specified.
func ValidateDIDUrl(didUrl string, method string, allowedNamespaces []string) error {
	did, path, query, fragment := TrySplitDIDUrl(didUrl)
	// Validate DID
	err := ValidateDID(did, method, allowedNamespaces)
	if err != nil {
		return err
	}
	// Validate path
	err = ValidatePath(path)
	if err != nil {
		return err
	}
	// Validate query
	err = ValidateQuery(query)
	if err != nil {
		return err
	}
	// Validate fragment
	err = ValidateFragment(fragment)
	if err != nil {
		return err
	}

	return nil
}

func ValidateFragment(fragment string) error {
	if !DIDFragmentRegexp.MatchString(fragment) {
		return ErrStaticDIDURLFragmentNotValid.Wrap(fragment)
	}
	return nil
}

func ValidateQuery(query string) error {
	if !DIDQueryRegexp.MatchString(query) {
		return ErrStaticDIDURLQueryNotValid.Wrap(query)
	}
	return nil
}

func ValidatePath(path string) error {
	if !DIDPathAbemptyRegexp.MatchString(path) {
		return ErrStaticDIDURLPathAbemptyNotValid.Wrap(path)
	}
	return nil
}

func IsValidDIDUrl(didUrl string, method string, allowedNamespaces []string) bool {
	err := ValidateDIDUrl(didUrl, method, allowedNamespaces)

	return nil == err
}
