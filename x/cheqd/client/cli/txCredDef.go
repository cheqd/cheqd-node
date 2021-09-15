package cli

import (
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateCredDef() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-credDef [schema_id] [tag] [signature_type] [value]",
		Short: "Creates a new credDef",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsSchema_id := string(args[0])
			argsTag := string(args[1])
			argsSignature_type := string(args[2])
			argsValue := string(args[3])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCredDef(string(argsSchema_id), string(argsTag), string(argsSignature_type), string(argsValue))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
