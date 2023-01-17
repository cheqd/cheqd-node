package integration

import (
	"fmt"
	"path/filepath"
	"strings"

	helpers "github.com/cheqd/cheqd-node/tests/integration/helpers"
)

func GetCaseName(path string) string {
	split := strings.Split(path, string(filepath.Separator))
	l := len(split)
	file := split[l-1]
	idiom := strings.Replace(file, ".json", "", 1)
	return fmt.Sprintf("%s %s %s", split[l-3], split[l-2], idiom)
}

func GetFile(path string) string {
	split := strings.Split(path, string(filepath.Separator))
	return split[len(split)-1]
}

func CreateTestJSON(dir string, content []byte) (string, error) {
	return helpers.WriteTmpFile(dir, content)
}
