package utils

import (
	"regexp"
	"strings"
)

var DidForbiddenSymbolsRegexp, _ = regexp.Compile("^[^#?&/\\\\]+$")

func SplitDidUrlIntoDidAndFragment(didUrl string) (string, string) {
	fragments := strings.Split(didUrl, "#")
	return fragments[0], fragments[1]
}

func ArrayContainsNotDid(array []string) (bool, int) {
	for i, did := range array {
		if IsNotDid(did) {
			return true, i
		}
	}

	return false, 0
}

func IsNotDid(did string) bool {
	return IsDid(did) == false
}

func IsDid(did string) bool {
	if len(did) == 0 {
		return false
	}

	if !DidForbiddenSymbolsRegexp.MatchString(did) {
		return false
	}

	fragments := strings.Split(did, ":")

	if len(fragments) <= 3 {
		return false
	}

	if fragments[0] != "did" {
		return false
	}

	if fragments[1] != MethodName {
		return false
	}

	if fragments[2] != MethodSpecificId {
		return false
	}

	return true
}
