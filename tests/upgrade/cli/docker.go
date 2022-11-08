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
	DOCKER_COMPOSE_ENV   = "docker-compose.env"
	DOCKER               = "docker"
	DOCKER_COMPOSE       = "compose"
	DOCKER_LOAD          = "load"
	DOCKER_IMAGE_NAME    = "cheqd-node-build.tar"
	DOCKER_IMAGE_ENV     = "BUILD_IMAGE"
	DOCKER_IMAGE_BUILD   = "cheqd/cheqd-node:build-latest"
	DOCKER_IMAGE_LATEST  = "ghcr.io/cheqd/cheqd-node:latest"
	DOCKER_HOME          = "/home/cheqd"
	DOCKER_USER          = "cheqd"
	DOCKER_USER_GROUP    = "cheqd"
	RUNNER_BIN_DIR       = "$(echo $RUNNER_BIN_DIR)"
	OPERATOR0            = "operator-0"
	OPERATOR1            = "operator-1"
	OPERATOR2            = "operator-2"
	OPERATOR3            = "operator-3"
	VALIDATOR0           = "validator-0"
	VALIDATOR1           = "validator-1"
	VALIDATOR2           = "validator-2"
	VALIDATOR3           = "validator-3"
	VALIDATORS           = 4
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

func LocalnetExecCopyAbsoluteWithPermissions(path string, destination string, container string) (string, error) {
	_, err := LocalnetExec("cp", path, filepath.Join(container+":"+destination))
	if err != nil {
		fmt.Println("Error copying file to container: ", err)
		return "", err
	}
	return LocalnetExecRestorePermissions(destination, container)
}

func LocalnetExecRestorePermissions(path string, container string) (string, error) {
	return LocalnetExec("exec", "-it", "--user", "root", container, "chown", "-R", DOCKER_USER+":"+DOCKER_USER_GROUP, path)
}

func LocalnetExecCopyKeys() (string, error) {
	for _, validator := range ValidatorNodes {
		_, err := LocalnetExecCopyKey(validator)
		if err != nil {
			return "", err
		}
	}
	return "", nil
}

func LocalnetExecCopyKey(validator string) (string, error) {
	_, err := LocalnetExec("cp", filepath.Join(DOCKER_LOCALNET_PATH, NETWORK_CONFIG_DIR, validator, KEYRING_DIR), filepath.Join(validator+":", DOCKER_HOME, ".cheqdnode"))
	if err != nil {
		return "", err
	}
	return LocalnetExec("exec", "-it", "--user", "root", validator, "chown", "-R", DOCKER_USER+":"+DOCKER_USER_GROUP, DOCKER_HOME)
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

func LocalnetExecUpWithBuildImage() (string, error) {
	return ExecWithEnv(
		[]string{DOCKER_IMAGE_ENV + "=" + DOCKER_IMAGE_BUILD},
		DOCKER, DOCKER_COMPOSE,
		"--env-file",
		filepath.Join(DOCKER_LOCALNET_PATH, DOCKER_COMPOSE_ENV),
		"-f",
		filepath.Join(DOCKER_LOCALNET_PATH, DOCKER_COMPOSE_FILE),
		"up",
		"--no-build",
	)
}

func LocalnetExecUpWithNewImage() (string, error) {
	err := SetNewDockerComposeEnv()
	if err != nil {
		return "", err
	}
	out, err := LocalnetLoadImage(DOCKER_LOAD_IMAGE_ARGS...)
	if err != nil {
		fmt.Println("Error on loading build image", out)
		return "", err
	}
	out, err = LocalnetExecUpWithBuildImage()
	fmt.Println("Rebooting localnet with new image", out)
	if err != nil {
		return "", err
	}
	return out, nil
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
		_, err := ExecDirect(RENAME_BINARY_CURRENT_TO_PREVIOUS_ARGS...)
		if err != nil {
			return "", err
		}
		_, err = ExecDirect(RENAME_BINARY_NEXT_TO_CURRENT_ARGS...)
		if err != nil {
			return "", err
		}
		return ExecDirect(RESTORE_BINARY_PERMISSIONS_ARGS...)
	case "next-to-previous":
		_, err := ExecDirect(RENAME_BINARY_CURRENT_TO_NEXT_ARGS...)
		if err != nil {
			return "", err
		}
		_, err = ExecDirect(RENAME_BINARY_PREVIOUS_TO_CURRENT_ARGS...)
		if err != nil {
			return "", err
		}
		return ExecDirect(RESTORE_BINARY_PERMISSIONS_ARGS...)
	default:
		return "", fmt.Errorf("invalid action")
	}
}
