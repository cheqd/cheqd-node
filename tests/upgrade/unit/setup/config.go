package setup

import "path/filepath"

const (
	GENERATED_JSON_DIR = "generated"
)

func JoinGenerated(path ...string) string {
	return filepath.Join(append([]string{GENERATED_JSON_DIR}, path...)...)
}
