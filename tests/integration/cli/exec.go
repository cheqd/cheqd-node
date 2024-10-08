package cli

import (
	"fmt"
	"os/exec"
	"time"

	"cosmossdk.io/errors"
)

func Exec(args ...string) (string, error) {
	cmd := exec.Command(CliBinaryName, args...)
	fmt.Println("cmd>>>>>>>>>>>>>>>>>>>>", cmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", errors.Wrap(err, string(out))
	}

	time.Sleep(2000 * time.Millisecond)

	return string(out), err
}
