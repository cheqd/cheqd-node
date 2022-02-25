package utils

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

func ValidateDID(did string, allowedNamespaces []string) error {
	// TODO: Implement
	return nil
}

func IsValidDID(did string, allowedNamespaces []string) bool {
	err := ValidateDID(did, allowedNamespaces)
	return err == nil
}

// SplitDID panics if did is not valid
func SplitDID(did string) (method string, namespace string, id string) {
	// TODO: Implement
	return "", "", ""
}

// DIDUrl: did:namespace:id[/path][?query][#fragment]
// TODO: Can path, query, fragment be set at the same time?
// TODO: Is service -> id URI or DIDUrl? What should we support?
// https://www.w3.org/TR/did-core/#did-url-syntax

func ValidateDIDUrl(didUrl string, allowedNamespaces []string) error {
	// TODO: Implement
	return nil
}

func IsValidDIDUrl(didUrl string, allowedNamespaces []string) bool {
	err := ValidateDIDUrl(didUrl, allowedNamespaces)
	return err == nil
}

// SplitDIDUrl panics if did is not valid
func SplitDIDUrl(didUrl string) (method string, namespace string, id string, path string, query string, fragment string) {
	// TODO: Implement
	return "", "", "", "", "", ""
}
