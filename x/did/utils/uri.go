package utils

import (
	"fmt"
	"regexp"
)

// Goes from RFC: https://www.rfc-editor.org/rfc/rfc3986#appendix-B
var ValidURIRegexp, _ = regexp.Compile(`^(([^:\/?#]+):)?(\/\/([^\/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`)

func ValidateURI(uri string) error {
	// Match with Regexp from RFC
	if !ValidURIRegexp.MatchString(uri) {
		return fmt.Errorf("URI: %s does not match regexp: %s", uri, ValidURIRegexp)
	}

	return nil
}
