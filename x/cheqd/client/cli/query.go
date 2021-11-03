package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd() *cobra.Command {
	// Group cheqd queries under a subcommand
	cmd := &cobra.Command{
		Use:                        v1.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", v1.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(CmdShowDid())

	return cmd
}
