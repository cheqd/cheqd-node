package cli

import (
	"os/exec"

	errorsmod "cosmossdk.io/errors"
)

func Exec(args ...string) (string, error) {
	cmd := exec.Command(CliBinaryName, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errorsmod.Wrap(err, string(out))
	}

	return string(out), err
}
