package helpers

import (
	"strings"
)

func TrimExtraLineOffset(input string, offset int) string {
	return strings.Join(strings.Split(input, "\n")[offset:], "")
}
