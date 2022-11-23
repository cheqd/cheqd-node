//go:build upgrade

package upgrade

import (
	"os"
	"path/filepath"
)

func RelGlob(relativePath ...string) ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	relativePathJoined := filepath.Join(relativePath...)
	fullPath := filepath.Join(cwd, relativePathJoined)

	return filepath.Glob(fullPath)
}
