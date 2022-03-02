package types

import (
	"github.com/cheqd/cheqd-node/x/cheqd/utils"
	"regexp"
)

var SplitDIDRegexp, _        = regexp.Compile(`did:([^:]+?)(:([^:]+?))?:([^:]+)$`)
var DidNamespaceRegexp, _    = regexp.Compile(`^[a-zA-Z0-9]$`)
// Base58 only allowed (without OolI and 0)
var UniqueIDRegexp, _        = regexp.Compile(`^[123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*$`)
// That for groups:
// Example: did:cheqd:testnet:fafdsffq11213343/path-to-s/ome-external-resource?query#key1???#123
// 1 - [^/?#]* - all the symbols except / and ? and # . This is the DID part                      (did:cheqd:testnet:fafdsffq11213343)
// 2 - [^?#]*  - all the symbols except ? and #. it means te section started from /, path-abempty (/path-to-s/ome-external-resource)
// 3 - \?([^#]*) - group for `query` part but with ? symbol 									  (?query)
// 4 - [^#]*     - group inside query string, match only exact query                              (query)
// 5 - #([^#]+[\$]?) - group for fragment, starts with #, includes #                              (#key1???)
// 6 - [^#]+[\$]?    - fragment only															  (key1???)
// {0,1} - means that number of fragments can be only 0 or 1
// Amount of query is not limited.
var SplitDIDURL, _           = regexp.Compile(`([^/?#]*)?([^?#]*)(\?([^#]*)){0,1}(#([^#]+$)){0,1}$`)
var DIDPathAbemptyRegexp, _  = regexp.Compile(`^[/a-zA-Z0-9\-\.\_\~\%\!\$\&\'\(\)\*\+\,\;\=\:\@]+$`)
var DIDQueryRegexp, _        = regexp.Compile(`^[/a-zA-Z0-9\-\.\_\~\%\!\$\&\'\(\)\*\+\,\;\=\:\@\/\?]*$`)
var DIDFragmentRegexp, _     = regexp.Compile(`^[/a-zA-Z0-9\-\.\_\~\%\!\$\&\'\(\)\*\+\,\;\=\:\@\/\?]*$`)

//[^/]+(/[^\?#]*)?(\?[^#]+)?(#.+)?
//Old implementation
//var DidForbiddenSymbolsRegexp, _ = regexp.Compile(`^[^#?&/\\]+$`)
//
//func SplitDidUrlIntoDidAndFragment(didUrl string) (string, string) {
//	fragments := strings.Split(didUrl, "#")
//	return fragments[0], fragments[1]
//}
//
//func IsDidFragment(prefix string, didUrl string) bool {
//	if !strings.Contains(didUrl, "#") {
//		return false
//	}
//
//	if didUrl[0] == '#' {
//		return true
//	}
//
//	did, _ := SplitDidUrlIntoDidAndFragment(didUrl)
//	return IsValidDid(prefix, did)
//}
//
//func IsFullDidFragment(prefix string, didUrl string) bool {
//	if !strings.Contains(didUrl, "#") {
//		return false
//	}
//
//	did, _ := SplitDidUrlIntoDidAndFragment(didUrl)
//	return IsValidDid(prefix, did)
//}
//
//func ResolveId(did string, methodId string) string {
//	result := methodId
//
//	methodDid, methodFragment := SplitDidUrlIntoDidAndFragment(methodId)
//	if len(methodDid) == 0 {
//		result = did + "#" + methodFragment
//	}
//
//	return result
//}
//
//func IsNotValidDIDArray(prefix string, array []string) (bool, int) {
//	for i, did := range array {
//		if !IsValidDid(prefix, did) {
//			return true, i
//		}
//	}
//
//	return false, 0
//}
//
//func IsNotValidDIDArrayFragment(prefix string, array []string) (bool, int) {
//	for i, did := range array {
//		if !IsDidFragment(prefix, did) {
//			return true, i
//		}
//	}
//
//	return false, 0
//}
//
//func IsValidDid(prefix string, did string) bool {
//	if len(did) == 0 {
//		return false
//	}
//
//	if !DidForbiddenSymbolsRegexp.MatchString(did) {
//		return false
//	}
//
//	// FIXME: Empty namespace must be allowed even if namespace is set in state
//	// https://github.com/cheqd/cheqd-node/blob/main/architecture/adr-list/adr-002-cheqd-did-method.md#method-specific-identifier
//	return strings.HasPrefix(did, prefix)
//}


// DID

func ValidateDID(did string, method string, allowedNamespaces []string) error {
	method, namespace, unique_id := SplitDID(did)

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

// SplitDID panics if did is not valid
func SplitDID(did string) (method string, namespace string, id string) {
	// Example: did:cheqd:testnet:base58str1ng1111
	// match [0] - the whole string
	// match [1] - cheqd                - method
	// match [2] - :testnet
	// match [3] - testnet              - namespace
	// match [4] - base58str1ng1111     - id
	match := SplitDIDRegexp.FindStringSubmatch(did)
	if len(match) > 0 {
		return match[1], match[3], match[4]
	}

	return "", "", ""
}

// SplitDIDUrl panics if did cannot be splitted properly
func SplitDIDUrl(didUrl string) (did string, path string , query string, fragment string) {
	match := SplitDIDURL.FindStringSubmatch(didUrl)
	return match[1], match[2], match[4], match[6]
}


// DIDUrl: did:namespace:id[/path][?query][#fragment]
// TODO: Can path, query, fragment be set at the same time?
// TODO: Is service -> id URI or DIDUrl? What should we support?
// https://www.w3.org/TR/did-core/#did-url-syntax

func ValidateDIDUrl(didUrl string, method string, allowedNamespaces []string) error {
	did, path, query, fragment := SplitDIDUrl(didUrl)
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