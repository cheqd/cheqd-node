package cli

import (
	"github.com/spf13/cobra"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/verim-id/verim-cosmos/x/verimcosmos/types"
)

func CmdCreateNym() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-nym [alais] [verkey] [did] [role]",
		Short: "Creates a new nym",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			argsAlais := string(args[0])
			argsVerkey := string(args[1])
			argsDid := string(args[2])
			argsRole := string(args[3])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgCreateNym(clientCtx.GetFromAddress().String(), string(argsAlais), string(argsVerkey), string(argsDid), string(argsRole))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdUpdateNym() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-nym [id] [alais] [verkey] [did] [role]",
		Short: "Update a nym",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			id, err := strconv.ParseUint(args[0], 10, 64)
			if err != nil {
				return err
			}

			argsAlais := string(args[1])
			argsVerkey := string(args[2])
			argsDid := string(args[3])
			argsRole := string(args[4])

			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			msg := types.NewMsgUpdateNym(clientCtx.GetFromAddress().String(), id, string(argsAlais), string(argsVerkey), string(argsDid), string(argsRole))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}

func CmdDeleteNym() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete-nym [id] [alais] [verkey] [did] [role]",
		Short: "Delete a nym by id",
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

			msg := types.NewMsgDeleteNym(clientCtx.GetFromAddress().String(), id)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}
			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
