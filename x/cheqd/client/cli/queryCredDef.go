package cli

import (
	"context"
	"github.com/cheqd/cheqd-node/x/cheqd/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdShowCredDef() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-cred-def [id]",
		Short: "shows a credDef",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)
			id := string(args[0])
			params := &types.QueryGetCredDefRequest{
				Id: id,
			}

			res, err := queryClient.CredDef(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
