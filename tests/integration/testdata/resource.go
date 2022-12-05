package testdata

import "github.com/cheqd/cheqd-node/tests/integration/helpers"

const (
	JSON_FILE_CONTENT = `{"test": "test"}`
)

func CreateTestJson(dir string) (string, error) {
	return helpers.WriteTmpFile(dir, []byte(JSON_FILE_CONTENT))
}
