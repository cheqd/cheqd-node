package utils

import (
	"errors"
	"fmt"
	"regexp"
)

// That for groups:
// Example: did:cheqd:testnet:fafdsffq11213343/path-to-s/ome-external-resource?query#key1???
// 1 - [^/?#]* - all the symbols except / and ? and # . This is the DID part                      (did:cheqd:testnet:fafdsffq11213343)
// 2 - [^?#]*  - all the symbols except ? and #. it means te section started from /, path-abempty (/path-to-s/ome-external-resource)
// 3 - \?([^#]*) - group for `query` part but with ? symbol 									  (?query)
// 4 - [^#]*     - group inside query string, match only exact query                              (query)
// 5 - #([^#]+[\$]?) - group for fragment, starts with #, includes #                              (#key1???)
// 6 - [^#]+[\$]?    - fragment only															  (key1???)
// Number of queries is not limited.
var SplitDIDURLRegexp, _ = regexp.Compile(`([^/?#]*)?([^?#]*)(\?([^#]*))?(#([^#]+$))?$`)

var (
	DIDPathAbemptyRegexp, _ = regexp.Compile(`^([/a-zA-Z0-9\-\.\_\~\!\$\&\'\(\)\*\+\,\;\=\:\@]*|(%[0-9A-Fa-f]{2})*)*$`)
	DIDQueryRegexp, _       = regexp.Compile(`^([/a-zA-Z0-9\-\.\_\~\!\$\&\'\(\)\*\+\,\;\=\:\@\/\?]*|(%[0-9A-Fa-f]{2})*)*$`)
	DIDFragmentRegexp, _    = regexp.Compile(`^([/a-zA-Z0-9\-\.\_\~\!\$\&\'\(\)\*\+\,\;\=\:\@\/\?]*|(%[0-9A-Fa-f]{2})*)*$`)
)

// TrySplitDIDUrl Validates generic format of DIDUrl. It doesn't validate path, query and fragment content.
// Call ValidateDIDUrl for further validation.
func TrySplitDIDUrl(didUrl string) (did string, path string, query string, fragment string, err error) {
	matches := SplitDIDURLRegexp.FindAllStringSubmatch(didUrl, -1)

	if len(matches) != 1 {
		return "", "", "", "", errors.New("unable to split did url into did, path, query and fragment")
	}

	match := matches[0]

	return match[1], match[2], match[4], match[6], nil
}

func MustSplitDIDUrl(didUrl string) (did string, path string, query string, fragment string) {
	did, path, query, fragment, err := TrySplitDIDUrl(didUrl)
	if err != nil {
		panic(err.Error())
	}
	return
}

func JoinDIDUrl(did string, path string, query string, fragment string) string {
	res := did + path

	if query != "" {
		res = res + "?" + query
	}

	if fragment != "" {
		res = res + "#" + fragment
	}

	return res
}

// ValidateDIDUrl checks method and allowed namespaces only when the corresponding parameters are specified.
func ValidateDIDUrl(didUrl string, method string, allowedNamespaces []string) error {
	did, path, query, fragment, err := TrySplitDIDUrl(didUrl)
	if err != nil {
		return err
	}

	// Validate DID
	err = ValidateDID(did, method, allowedNamespaces)
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
		return fmt.Errorf("did url fragmnt must match the following regexp: %s", DIDFragmentRegexp)
	}
	return nil
}

func ValidateQuery(query string) error {
	if !DIDQueryRegexp.MatchString(query) {
		return fmt.Errorf("did url query must match the following regexp: %s", DIDQueryRegexp)
	}
	return nil
}

func ValidatePath(path string) error {
	if !DIDPathAbemptyRegexp.MatchString(path) {
		return fmt.Errorf("did url path abempty must match the following regexp: %s", DIDPathAbemptyRegexp)
	}
	return nil
}

func IsValidDIDUrl(didUrl string, method string, allowedNamespaces []string) bool {
	err := ValidateDIDUrl(didUrl, method, allowedNamespaces)

	return nil == err
}
