package cli

import (
	"encoding/base64"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateDid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-did [did]",
		Short: "Creates a new did",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsDidBase64 := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var msg types.MsgWriteRequest
			data, err := base64.StdEncoding.DecodeString(argsDidBase64)
			if err != nil {
				return err
			}

			if err := msg.Unmarshal(data); err != nil {
				return err
			}

			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateDid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-did [did]",
		Short: "Update a did",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsDidBase64 := args[0]

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			var msg types.MsgWriteRequest
			data, err := base64.StdEncoding.DecodeString(argsDidBase64)
			if err != nil {
				return err
			}

			if err := msg.Unmarshal(data); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
