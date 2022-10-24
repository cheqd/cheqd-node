package cli

import (
	integrationcli "github.com/cheqd/cheqd-node/tests/integration/cli"
)

func Exec(args ...string) (string, error) {
	return integrationcli.Exec(args...)
}