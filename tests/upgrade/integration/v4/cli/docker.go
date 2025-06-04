package cli

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	errorsmod "cosmossdk.io/errors"
	"github.com/joho/godotenv"
)

const (
	DockerLocalnetPath = "../../../../docker/localnet"
	DockerComposeFile  = "docker-compose.yml"
	DockerComposeEnvML = "mainnet-latest.env"
	DockerComposeEnvBL = "build-latest.env"
	Docker             = "docker"
	DockerCompose      = "compose"
	DockerHome         = "/home/cheqd"
	DockerUser         = "cheqd"
	DockerUserGroup    = "cheqd"
	Operator0          = "operator-0"
	Operator1          = "operator-1"
	Operator2          = "operator-2"
	Operator3          = "operator-3"
	Validator0         = "validator-0"
	Validator1         = "validator-1"
	Validator2         = "validator-2"
	Validator3         = "validator-3"
	ValidatorsCount    = 4
)

type OperatorAccountType map[string]string

var OperatorAccounts = OperatorAccountType{
	Validator0: Operator0,
	Validator1: Operator1,
	Validator2: Operator2,
	Validator3: Operator3,
}

var ValidatorNodes = []string{Validator0, Validator1, Validator2, Validator3}

var (
	DockerComposeLatestArgs = []string{
		"-f", filepath.Join(DockerLocalnetPath, DockerComposeFile),
		"--env-file", filepath.Join(DockerLocalnetPath, DockerComposeEnvML),
	}
	DockerComposeBuildArgs = []string{
		"-f", filepath.Join(DockerLocalnetPath, DockerComposeFile),
		"--env-file", filepath.Join(DockerLocalnetPath, DockerComposeEnvBL),
	}
)

func LocalnetExec(envArgs []string, args ...string) (string, error) {
	args = append(append([]string{DockerCompose}, envArgs...), args...)
	cmd := exec.Command(Docker, args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return string(out), errorsmod.Wrap(err, string(out))
	}
	return string(out), err
}

func LocalnetExecExec(container string, args ...string) (string, error) {
	args = append([]string{"exec", container}, args...)
	return LocalnetExec(DockerComposeLatestArgs, args...)
}

func LocalnetExecUp() (string, error) {
	return LocalnetExec(DockerComposeLatestArgs, "up", "--detach", "--no-build")
}

func LocalnetExecDown() (string, error) {
	return LocalnetExec(DockerComposeLatestArgs, "down")
}

func LocalnetExecCopyAbsoluteWithPermissions(path string, destination string, container string) (string, error) {
	_, err := LocalnetExec(DockerComposeLatestArgs, "cp", path, container+":"+destination)
	if err != nil {
		fmt.Println("Error copying file to container: ", err)
		return "", err
	}
	return LocalnetExecRestorePermissions(destination, container)
}

func LocalnetExecRestorePermissions(path string, container string) (string, error) {
	return LocalnetExec(DockerComposeLatestArgs, "exec", "-it", "--user", "root", container, "chown", "-R", DockerUser+":"+DockerUserGroup, path)
}

func LocalnetExecRunWithVolume(containerId string, command []string) (string, error) {
	// Load .env file
	envFile := filepath.Join(DockerLocalnetPath, DockerComposeEnvBL)
	err := godotenv.Load(envFile)
	if err != nil {
		return "", fmt.Errorf("error loading %s file:%v", envFile, err)
	}

	// Get variables
	imageName := os.Getenv("BUILD_IMAGE")
	if imageName == "" {
		return "", fmt.Errorf("BUILD_IMAGE is not set in .env")
	}

	return ExecDirect(append([]string{Docker, "run", "--rm", "--volumes-from", containerId, imageName}, command...)...)
}

func LocalnetStartContainer(container string) (string, error) {
	containerId, err := GetContainerIDByName(container)
	if err != nil {
		return "", err
	}
	return ExecDirect(Docker, "start", containerId)
}

func LocalnetStopContainerWithId(containerId string) (string, error) {
	return ExecDirect(Docker, "stop", containerId)
}

func GetContainerIDByName(container string) (string, error) {
	out, err := ExecDirect(Docker, "ps", "-q", "-a", "-f", "name="+container)
	if err != nil {
		return "", fmt.Errorf("failed to get container ID: %v", err)
	}

	containerId := strings.TrimSpace(out)
	if containerId == "" {
		log.Fatalf("no container id found for: %s", container)
	}

	return containerId, nil
}
