package cli

import (
	"os/exec"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

func Exec(args ...string) (string, error) {
	cmd := exec.Command(CLI_BINARY_NAME, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, string(out))
	}

	return string(out), err
}
