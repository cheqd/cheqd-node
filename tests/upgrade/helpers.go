//go:build upgrade

package upgrade

import (
	"fmt"
	"path/filepath"
	"strings"
)

func GetCase(path string) (case_ string, file string) {
	split := strings.Split(path, string(filepath.Separator))
	l := len(split)
	file = split[l-1]
	idiom := strings.Replace(file, ".json", "", 1)
	return fmt.Sprintf("%s %s %s", split[l-3], split[l-2], idiom), file
}

func GetFile(path string) string {
	split := strings.Split(path, string(filepath.Separator))
	return split[len(split)-1]
}
