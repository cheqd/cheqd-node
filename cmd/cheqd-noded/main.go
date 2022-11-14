package main

import (
	"os"

	"github.com/canow-co/cheqd-node/app"
	"github.com/canow-co/cheqd-node/cmd/cheqd-noded/cmd"
	svrcmd "github.com/cosmos/cosmos-sdk/server/cmd"
)

func main() {
	rootCmd, _ := cmd.NewRootCmd()
	if err := svrcmd.Execute(rootCmd, app.DefaultNodeHome); err != nil {
		os.Exit(1)
	}
}
