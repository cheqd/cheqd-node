package cli

import (
	"context"
	"github.com/cheqd/cheqd-node/x/cheqd/types/v1"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdShowDid() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "show-did [id]",
		Short: "shows a did",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := v1.NewQueryClient(clientCtx)

			id := args[0]
			params := &v1.QueryGetDidRequest{
				Id: id,
			}

			res, err := queryClient.Did(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(res)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
