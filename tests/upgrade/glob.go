//go:build upgrade

package upgrade

import (
	"os"
	"path/filepath"
)

func Glob(path string) ([]string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return filepath.Glob(filepath.Join(cwd, path))
}
