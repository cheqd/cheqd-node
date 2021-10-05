package cli

import (
	"encoding/json"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateCredDef() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-cred-def [id] [schema_id] [tag] [signature_type] [value]",
		Short: "Creates a new credDef",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsId := args[0]
			argsSchemaId := args[1]
			argsTag := args[2]
			argsSignatureType := args[3]
			argsValue := args[4]

			var value types.MsgCreateCredDef_ClType
			if err := json.Unmarshal([]byte(argsValue), &value); err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateCredDef(argsId, argsSchemaId, argsTag, argsSignatureType, &value)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
