package cli

import (
	"fmt"
	"os/exec"
	"path/filepath"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DOCKER_LOCALNET_PATH  = "../../docker/localnet"
	DOCKER_COMPOSE_FILE   = "docker-compose.yml"
	DOCKER_COMPOSE_ENV_ML = "mainnet-latest.env"
	DOCKER_COMPOSE_ENV_BL = "build-latest.env"
	DOCKER                = "docker"
	DOCKER_COMPOSE        = "compose"
	DOCKER_HOME           = "/home/cheqd"
	DOCKER_USER           = "cheqd"
	DOCKER_USER_GROUP     = "cheqd"
	OPERATOR0             = "operator-0"
	OPERATOR1             = "operator-1"
	OPERATOR2             = "operator-2"
	OPERATOR3             = "operator-3"
	VALIDATOR0            = "validator-0"
	VALIDATOR1            = "validator-1"
	VALIDATOR2            = "validator-2"
	VALIDATOR3            = "validator-3"
	VALIDATORS            = 4
)

type OperatorAccount map[string]string

var OperatorAccounts OperatorAccount = OperatorAccount{
	VALIDATOR0: OPERATOR0,
	VALIDATOR1: OPERATOR1,
	VALIDATOR2: OPERATOR2,
	VALIDATOR3: OPERATOR3,
}

var ValidatorNodes = []string{VALIDATOR0, VALIDATOR1, VALIDATOR2, VALIDATOR3}

var (
	DOCKER_COMPOSE_LATEST_ARGS = []string{
		"-f", filepath.Join(DOCKER_LOCALNET_PATH, DOCKER_COMPOSE_FILE),
		"--env-file", filepath.Join(DOCKER_LOCALNET_PATH, DOCKER_COMPOSE_ENV_ML),
	}
	DOCKER_COMPOSE_BUILD_ARGS = []string{
		"-f", filepath.Join(DOCKER_LOCALNET_PATH, DOCKER_COMPOSE_FILE),
		"--env-file", filepath.Join(DOCKER_LOCALNET_PATH, DOCKER_COMPOSE_ENV_BL),
	}
)

func LocalnetExec(envArgs []string, args ...string) (string, error) {
	args = append(append([]string{DOCKER_COMPOSE}, envArgs...), args...)
	cmd := exec.Command(DOCKER, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", sdkerrors.Wrap(err, string(out))
	}
	return string(out), err
}

func LocalnetExecExec(container string, args ...string) (string, error) {
	args = append([]string{"exec", container}, args...)
	return LocalnetExec(DOCKER_COMPOSE_LATEST_ARGS, args...)
}

func LocalnetExecUp() (string, error) {
	return LocalnetExec(DOCKER_COMPOSE_LATEST_ARGS, "up", "--detach", "--no-build")
}

func LocalnetExecDown() (string, error) {
	return LocalnetExec(DOCKER_COMPOSE_LATEST_ARGS, "down")
}

func LocalnetExecCopyAbsoluteWithPermissions(path string, destination string, container string) (string, error) {
	_, err := LocalnetExec(DOCKER_COMPOSE_LATEST_ARGS, "cp", path, filepath.Join(container+":"+destination))
	if err != nil {
		fmt.Println("Error copying file to container: ", err)
		return "", err
	}
	return LocalnetExecRestorePermissions(destination, container)
}

func LocalnetExecRestorePermissions(path string, container string) (string, error) {
	return LocalnetExec(DOCKER_COMPOSE_LATEST_ARGS, "exec", "-it", "--user", "root", container, "chown", "-R", DOCKER_USER+":"+DOCKER_USER_GROUP, path)
}
