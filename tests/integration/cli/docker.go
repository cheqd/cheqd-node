package cli

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	errorsmod "cosmossdk.io/errors"
)

const (
	DockerLocalnetPath = "../../docker/localnet"
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
	// Check if args contain "--output json"
	containsJSONOutput := false
	for i := 0; i < len(args)-1; i++ {
		if args[i] == "--output" && args[i+1] == "json" {
			containsJSONOutput = true
			break
		}
	}
	fmt.Printf("containsJSONOutput: %v\n", containsJSONOutput)

	// Combine Docker command with envArgs and args
	args = append(append([]string{DockerCompose}, envArgs...), args...)
	cmd := exec.Command(Docker, args...)
	out, err := cmd.CombinedOutput()

	// If "--output json" is present, attempt to extract JSON
	if containsJSONOutput {
		fmt.Printf("\"extracting json\": %v\n", "extracting json")
		extractedJSON, err := extractOnlyJSON(string(out))
		if err != nil {
			return "", errorsmod.Wrap(err, string(out))
		}
		return extractedJSON, nil
	}

	// If no "--output json", return the raw output
	if err != nil {
		return string(out), errorsmod.Wrap(err, string(out))
	}
	return string(out), nil
}

func extractOnlyJSON(out string) (string, error) {
	fmt.Printf("out: %v\n", out)
	// Find the first '{' and last '}' to extract the JSON portion
	jsonStart := strings.Index(out, "{")
	jsonEnd := strings.LastIndex(out, "}")
	if jsonStart == -1 || jsonEnd == -1 {
		return "", fmt.Errorf("no valid JSON found in output")
	}

	return out[jsonStart : jsonEnd+1], nil
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
