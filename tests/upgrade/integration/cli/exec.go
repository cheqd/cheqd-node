package cli

import (
	"os"
	"os/exec"

	errorsmod "cosmossdk.io/errors"
	integrationcli "github.com/cheqd/cheqd-node/tests/integration/cli"
)

func Exec(args ...string) (string, error) {
	return integrationcli.Exec(args...)
}

func ExecDirect(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errorsmod.Wrap(err, string(out))
	}

	return string(out), err
}

func ExecWithEnv(env []string, args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errorsmod.Wrap(err, string(out))
	}

	return string(out), err
}
