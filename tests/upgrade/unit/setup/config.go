package setup

import "path/filepath"

const (
	GeneratedJSONDir = "generated"
)

func JoinGenerated(path ...string) string {
	return filepath.Join(append([]string{GeneratedJSONDir}, path...)...)
}
