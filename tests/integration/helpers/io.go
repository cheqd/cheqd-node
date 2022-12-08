package helpers

import (
	"io/fs"
	"os"
	"path"

	"github.com/google/uuid"
)

func WriteTmpFile(tmpDir string, content []byte) (string, error) {
	name := uuid.NewString()
	file := path.Join(tmpDir, name)
	err := os.WriteFile(file, []byte(content), fs.ModePerm)
	if err != nil {
		return "", err
	}
	return file, nil
}

func MustWriteTmpFile(tmpDir string, content []byte) string {
	file, err := WriteTmpFile(tmpDir, content)
	if err != nil {
		panic(err)
	}
	return file
}
