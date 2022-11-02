package cli

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	ROOT_REL_PATH        = "../.."
	DOCKER_LOCALNET      = "localnet"
	DOCKER_LOCALNET_PATH = "../../docker/localnet"
	DOCKER_COMPOSE_FILE  = "docker-compose.yml"
	DOCKER               = "docker"
	DOCKER_COMPOSE       = "compose"
	DOCKER_LOAD          = "load"
	DOCKER_IMAGE_NAME    = "cheqd-node-image.tar"
	RUNNER_BIN_DIR       = "$(echo $RUNNER_BIN_DIR)"
	OPERATOR0            = "operator0"
	OPERATOR1            = "operator1"
	OPERATOR2            = "operator2"
	OPERATOR3            = "operator3"
	VALIDATOR0           = "validator-0"
	VALIDATOR1           = "validator-1"
	VALIDATOR2           = "validator-2"
	VALIDATOR3           = "validator-3"
)

type OperatorAccount map[string]string

var OperatorAccounts OperatorAccount = OperatorAccount{
	VALIDATOR0: OPERATOR0,
	VALIDATOR1: OPERATOR1,
	VALIDATOR2: OPERATOR2,
	VALIDATOR3: OPERATOR3,
}

var (
	DOCKER_COMPOSE_ARGS = []string{
		"-f", filepath.Join(DOCKER_LOCALNET_PATH, DOCKER_COMPOSE_FILE),
	}
	DOCKER_LOAD_IMAGE_ARGS = []string{
		"-i", filepath.Join(ROOT_REL_PATH, DOCKER_IMAGE_NAME),
	}
	RENAME_BINARY_CURRENT_TO_PREVIOUS_ARGS = []string{
		"mv",
		filepath.Join(RUNNER_BIN_DIR, CLI_BINARY_NAME),
		filepath.Join(RUNNER_BIN_DIR, CLI_BINARY_NAME_PREVIOUS),
	}
	RENAME_BINARY_NEXT_TO_CURRENT_ARGS = []string{
		"mv",
		filepath.Join(RUNNER_BIN_DIR, CLI_BINARY_NAME_NEXT),
		filepath.Join(RUNNER_BIN_DIR, CLI_BINARY_NAME),
	}
	RENAME_BINARY_PREVIOUS_TO_CURRENT_ARGS = []string{
		"mv",
		filepath.Join(RUNNER_BIN_DIR, CLI_BINARY_NAME_PREVIOUS),
		filepath.Join(RUNNER_BIN_DIR, CLI_BINARY_NAME),
	}
	RENAME_BINARY_CURRENT_TO_NEXT_ARGS = []string{
		"mv",
		filepath.Join(RUNNER_BIN_DIR, CLI_BINARY_NAME),
		filepath.Join(RUNNER_BIN_DIR, CLI_BINARY_NAME_NEXT),
	}
	RESTORE_BINARY_PERMISSIONS_ARGS = []string{
		"sudo",
		"chmod",
		"-x",
		filepath.Join(RUNNER_BIN_DIR, CLI_BINARY_NAME),
	}
)

func LocalnetExec(args ...string) (string, error) {
	args = append(append([]string{DOCKER_COMPOSE}, DOCKER_COMPOSE_ARGS...), args...)
	cmd := exec.Command(DOCKER, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", sdkerrors.Wrap(err, string(out))
	}
	return string(out), err
}

func LocalnetExecExec(container string, args ...string) (string, error) {
	args = append([]string{"exec", container}, args...)
	return LocalnetExec(args...)
}

func LocalnetExecUp() (string, error) {
	return LocalnetExec("up", "--detach", "--no-build")
}

func LocalnetExecDown() (string, error) {
	return LocalnetExec("down")
}

func LocalnetLoadImage(args ...string) (string, error) {
	args = append([]string{DOCKER_LOAD}, args...)
	cmd := exec.Command(DOCKER, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", sdkerrors.Wrap(err, string(out))
	}
	return string(out), err
}

func LocalnetExecUpWithNewImage() (string, error) {
	err := SetNewDockerComposeEnv()
	if err != nil {
		return "", err
	}
	_, err = LocalnetLoadImage(DOCKER_LOAD_IMAGE_ARGS...)
	if err != nil {
		return "", err
	}
	return LocalnetExecUp()
}

func SetOldDockerComposeEnv() error {
	os.Setenv("CHEQD_IMAGE_FROM", CHEQD_IMAGE_FROM)
	os.Setenv("CHEQD_TAG_FROM", CHEQD_TAG_FROM)
	return nil
}

func SetNewDockerComposeEnv() error {
	os.Setenv("CHEQD_IMAGE_TO", CHEQD_IMAGE_TO)
	os.Setenv("CHEQD_TAG_TO", CHEQD_TAG_TO)
	return nil
}

func ReplaceBinaryWithPermissions(action string) (string, error) {
	switch action {
	case "previous-to-next":
		_, err := Exec(RENAME_BINARY_CURRENT_TO_PREVIOUS_ARGS...)
		if err != nil {
			return "", err
		}
		_, err = Exec(RENAME_BINARY_NEXT_TO_CURRENT_ARGS...)
		if err != nil {
			return "", err
		}
		return Exec(RESTORE_BINARY_PERMISSIONS_ARGS...)
	case "next-to-previous":
		_, err := Exec(RENAME_BINARY_CURRENT_TO_NEXT_ARGS...)
		if err != nil {
			return "", err
		}
		_, err = Exec(RENAME_BINARY_PREVIOUS_TO_CURRENT_ARGS...)
		if err != nil {
			return "", err
		}
		return Exec(RESTORE_BINARY_PERMISSIONS_ARGS...)
	default:
		return "", fmt.Errorf("invalid action")
	}
}
