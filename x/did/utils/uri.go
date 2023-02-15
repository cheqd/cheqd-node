package utils

import (
	"fmt"
	"regexp"
)

// ValidURIRegexp ...
// Goes from RFC: https://www.rfc-editor.org/rfc/rfc3986#appendix-B
var ValidURIRegexp = regexp.MustCompile(`^(([^:\/?#]+):)?(\/\/([^\/?#]*))?([^?#]*)(\?([^#]*))?(#(.*))?`)

func ValidateURI(uri string) error {
	// Match with Regexp from RFC
	if !ValidURIRegexp.MatchString(uri) {
		return fmt.Errorf("URI: %s does not match regexp: %s", uri, ValidURIRegexp)
	}

	return nil
}
