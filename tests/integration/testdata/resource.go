package testdata

import (
	"encoding/base64"

	"github.com/cheqd/cheqd-node/tests/integration/helpers"
)

const (
	DEFAULT_FILE_CONTENT = `<p>Test file content</p>` // anything valid but image, json is fine
	JSON_FILE_CONTENT    = `{"test": "test"}`
	IMAGE_FILE_CONTENT   = `iVBORw0KGgoAAAANSUhEUgAAAQAAAAEAAQMAAABmvDolAAAAA1BMVEW10NBjBBbqAAAAH0lEQVRoge3BAQ0AAADCoPdPbQ43oAAAAAAAAAAAvg0hAAABmmDh1QAAAABJRU5ErkJggg==`
)

func CreateTestJson(dir string) (string, error) {
	return helpers.WriteTmpFile(dir, []byte(JSON_FILE_CONTENT))
}

func CreateTestImage(dir string) (string, error) {
	png, err := base64.StdEncoding.DecodeString(IMAGE_FILE_CONTENT)
	if err != nil {
		return "", err
	}
	return helpers.WriteTmpFile(dir, png)
}

func CreateTestDefault(dir string) (string, error) {
	return helpers.WriteTmpFile(dir, []byte(DEFAULT_FILE_CONTENT))
}
