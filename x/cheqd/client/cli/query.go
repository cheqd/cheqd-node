package cli

import (
	"fmt"
	// "strings"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	// sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

// GetQueryCmd returns the cli query commands for this module
func GetQueryCmd(queryRoute string) *cobra.Command {
	// Group cheqd queries under a subcommand
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("Querying commands for the %s module", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1

	cmd.AddCommand(CmdListCred_def())
	cmd.AddCommand(CmdShowCred_def())

	cmd.AddCommand(CmdListSchema())
	cmd.AddCommand(CmdShowSchema())

	cmd.AddCommand(CmdListAttrib())
	cmd.AddCommand(CmdShowAttrib())

	cmd.AddCommand(CmdListDid())
	cmd.AddCommand(CmdShowDid())

	cmd.AddCommand(CmdListNym())
	cmd.AddCommand(CmdShowNym())

	return cmd
}
