package utils

import (
	"regexp"
	"strings"
)

var DidForbiddenSymbolsRegexp, _ = regexp.Compile(`^[^#?&/\\]+$`)
// Base58 allowed symbols:
//                                                          "123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz"
// Regexp from https://regex101.com/library/iT8tI7
var DidAllowedSymbolsRegexp, _ = regexp.Compile(`^[1]([123456789ABCDEFGHJKLMNPQRSTUVWXYZabcdefghijkmnopqrstuvwxyz]*)$`)

func SplitDidUrlIntoDidAndFragment(didUrl string) (string, string) {
	fragments := strings.Split(didUrl, "#")
	return fragments[0], fragments[1]
}

func IsDidFragment(prefix string, didUrl string) bool {
	if !strings.Contains(didUrl, "#") {
		return false
	}

	if didUrl[0] == '#' {
		return true
	}

	did, _ := SplitDidUrlIntoDidAndFragment(didUrl)
	return IsValidDid(prefix, did)
}

func IsFullDidFragment(prefix string, didUrl string) bool {
	if !strings.Contains(didUrl, "#") {
		return false
	}

	did, _ := SplitDidUrlIntoDidAndFragment(didUrl)
	return IsValidDid(prefix, did)
}

func ResolveId(did string, methodId string) string {
	result := methodId

	methodDid, methodFragment := SplitDidUrlIntoDidAndFragment(methodId)
	if len(methodDid) == 0 {
		result = did + "#" + methodFragment
	}

	return result
}

func IsNotValidDIDArray(prefix string, array []string) (bool, int) {
	for i, did := range array {
		if !IsValidDid(prefix, did) {
			return true, i
		}
	}

	return false, 0
}

func IsNotValidDIDArrayFragment(prefix string, array []string) (bool, int) {
	for i, did := range array {
		if !IsDidFragment(prefix, did) {
			return true, i
		}
	}

	return false, 0
}

func IsValidDid(prefix string, did string) bool {

	if len(did) == 0 {
		return false
	}

	if !DidForbiddenSymbolsRegexp.MatchString(did) {
		return false
	}

	// Checks for unique identifier

	unique_id := strings.TrimPrefix(did, prefix)

	if len(unique_id) != 16 && len(unique_id) != 32 {
		return false
	}

	if !DidAllowedSymbolsRegexp.MatchString(unique_id) {
		return false
	}

	// FIXME: Empty namespace must be allowed even if namespace is set in state
	// https://github.com/cheqd/cheqd-node/blob/main/architecture/adr-list/adr-002-cheqd-did-method.md#method-specific-identifier
	return strings.HasPrefix(did, prefix)
}
