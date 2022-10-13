package testdata

import (
	"io/fs"
	"os"
	"path"
)

const (
	JSON_FILE_NAME    = "test.json"
	JSON_FILE_CONTENT = `{"test": "test"}`
)

func CreateTestFile(dir string, name string, content []byte) (string, error) {
	file := path.Join(dir, name)
	err := os.WriteFile(file, []byte(content), fs.ModePerm)
	if err != nil {
		return "", err
	}
	return file, nil
}

func CreateTestJson(dir string) (string, error) {
	return CreateTestFile(dir, JSON_FILE_NAME, []byte(JSON_FILE_CONTENT))
}
