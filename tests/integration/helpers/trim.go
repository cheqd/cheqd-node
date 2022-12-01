package helpers

import (
	"strings"
)

func TrimExtraLineOffset(input string, offset int) string {
	return strings.Join(strings.Split(input, "\n")[offset:], "")
}

func TrimImportedStdout(lines []string) (output string, trimmed []string) {
	trimmed = make([]string, len(lines))
	for i, line := range lines {
		if !strings.Contains(line, "Successfully migrated key") {
			trimmed[i] = line
		}
	}

	return strings.Join(trimmed, "\n"), trimmed
}
