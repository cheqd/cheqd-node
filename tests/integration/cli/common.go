package cli

import (
	"os/exec"

	"github.com/cosmos/cosmos-sdk/types/errors"
)

const CLI_BINARY_NAME = "cheqd-noded"

const (
	CHAIN_ID        = "cheqd"
	KEYRING_BACKEND = "test"
	OUTPUT_FORMAT   = "json"
	GAS             = "auto"
	GAS_ADJUSTMENT  = "2.0"
	GAS_PRICES      = "25ncheq"
)

func Exec(args ...string) (string, error) {
	cmd := exec.Command(CLI_BINARY_NAME, args...)
	out, err := cmd.CombinedOutput()

	println(string(out))

	if err != nil {
		return "", errors.Wrap(err, string(out))
	}

	return string(out), err
}
