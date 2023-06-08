package cli

import (
	"os/exec"

	sdkerrors "cosmossdk.io/errors"
)

func Exec(args ...string) (string, error) {
	cmd := exec.Command(CliBinaryName, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", sdkerrors.Wrap(err, string(out))
	}

	return string(out), err
}
