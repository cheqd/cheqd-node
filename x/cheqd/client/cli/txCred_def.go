package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
)

func CmdCreateCred_def() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-cred_def [schema_id] [tag] [signature_type] [value]",
		Short: "Creates a new cred_def",
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

			msg := types.NewMsgCreateCred_def(clientCtx.GetFromAddress().String(), string(argsSchema_id), string(argsTag), string(argsSignature_type), string(argsValue))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateCred_def() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-cred_def [id] [schema_id] [tag] [signature_type] [value]",
		Short: "Update a cred_def",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argsSchema_id := string(args[1])
			argsTag := string(args[2])
			argsSignature_type := string(args[3])
			argsValue := string(args[4])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateCred_def(clientCtx.GetFromAddress().String(), id, string(argsSchema_id), string(argsTag), string(argsSignature_type), string(argsValue))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteCred_def() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-cred_def [id] [schema_id] [tag] [signature_type] [value]",
		Short: "Delete a cred_def by id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgDeleteCred_def(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
