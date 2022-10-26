package cli

import (
	"os"
	"os/exec"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	DOCKER_LOCALNET         = "localnet"
	DOCKER_LOCALNET_PATH    = "../../../docker/localnet"
	DOCKER_COMPOSE          = "docker compose"
	DOCKER_IN_LOCALNET_PATH = DOCKER_LOCALNET_PATH + "/" + DOCKER_COMPOSE
	OPERATOR0               = "operator0"
	OPERATOR1               = "operator1"
	OPERATOR2               = "operator2"
	OPERATOR3               = "operator3"
	VALIDATOR0              = "validator-0"
	VALIDATOR1              = "validator-1"
	VALIDATOR2              = "validator-2"
	VALIDATOR3              = "validator-3"
)

type OperatorAccount map[string]string

var OperatorAccounts OperatorAccount = OperatorAccount{
	VALIDATOR0: OPERATOR0,
	VALIDATOR1: OPERATOR1,
	VALIDATOR2: OPERATOR2,
	VALIDATOR3: OPERATOR3,
}

func LocalnetExec(args ...string) (string, error) {
	cmd := exec.Command(DOCKER_IN_LOCALNET_PATH, args...)
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
	return LocalnetExec("up", "-d")
}

func LocalnetExecDown() (string, error) {
	return LocalnetExec("down")
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
