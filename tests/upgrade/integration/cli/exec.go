package cli

import (
	"os"
	"fmt"
	"strings"
	"os/exec"
	"path/filepath"

	integrationcli "github.com/cheqd/cheqd-node/tests/integration/cli"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

func Exec(args ...string) (string, error) {
	return integrationcli.Exec(args...)
}

func ExecDirect(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	if os.Getenv("DEBUG") == "true" {
		fmt.Println("DEBUG: Command for run: ", strings.Join(cmd.Args, " "))
	}
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, string(out))
	}

	return string(out), err
}

func ExecDirectWithHome(operator string, args ...string) (string, error) {
	var node_home = os.Getenv("CHEQD_NODE_HOME")
	if node_home == "" {
		node_home = filepath.Join("../../../docker/localnet", NETWORK_CONFIG_DIR, operator)
	}
	return ExecDirect(append(args, "--home", node_home)...)
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
