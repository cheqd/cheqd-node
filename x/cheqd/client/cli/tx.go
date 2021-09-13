package cli

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	// "github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
)

// GetTxCmd returns the transaction commands for this module
func GetTxCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      fmt.Sprintf("%s transactions subcommands", types.ModuleName),
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	// this line is used by starport scaffolding # 1

	cmd.AddCommand(CmdCreateCred_def())
	cmd.AddCommand(CmdUpdateCred_def())
	cmd.AddCommand(CmdDeleteCred_def())

	cmd.AddCommand(CmdCreateSchema())
	cmd.AddCommand(CmdUpdateSchema())
	cmd.AddCommand(CmdDeleteSchema())

	cmd.AddCommand(CmdCreateAttrib())
	cmd.AddCommand(CmdUpdateAttrib())
	cmd.AddCommand(CmdDeleteAttrib())

	cmd.AddCommand(CmdCreateDid())
	cmd.AddCommand(CmdUpdateDid())
	cmd.AddCommand(CmdDeleteDid())

	cmd.AddCommand(CmdCreateNym())
	cmd.AddCommand(CmdUpdateNym())
	cmd.AddCommand(CmdDeleteNym())

	return cmd
}
