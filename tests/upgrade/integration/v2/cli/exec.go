package cli

import (
	"os"
	"os/exec"

	integrationcli "github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func Exec(args ...string) (string, error) {
	return integrationcli.Exec(args...)
}

func ExecDirect(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, string(out))
	}

	return string(out), err
}

func ExecWithEnv(env []string, args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = append(os.Environ(), env...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, string(out))
	}

	return string(out), err
}
