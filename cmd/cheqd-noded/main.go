package main

import (
	"fmt"
	"os"

	"github.com/cheqd/cheqd-node/app"
	"github.com/cheqd/cheqd-node/cmd/cheqd-noded/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.Name, app.DefaultNodeHome); err != nil {
		fmt.Fprintln(rootCmd.OutOrStderr(), err)
		os.Exit(1)
	}
}
