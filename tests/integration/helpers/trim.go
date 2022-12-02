package helpers

import (
	"strings"
)

func TrimExtraLineOffset(input string, offset int) string {
	return strings.Join(strings.Split(input, "\n")[offset:], "")
}

func TrimImportedStdout(output string) string {
	lines := strings.Split(output, "\n")
	trimmed := make([]string, len(lines))
	for i, line := range lines {
		if !strings.Contains(line, "Successfully migrated key") && !strings.Contains(line, "gas estimate:") {
			trimmed[i] = line
		}
	}

	return strings.Join(trimmed, "\n")
}
