package cli

import (
	"os/exec"

	"cosmossdk.io/errors"
)

func Exec(args ...string) (string, error) {
	cmd := exec.Command(CliBinaryName, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, string(out))
	}

	// time.Sleep(10000 * time.Millisecond)

	return string(out), err
}
