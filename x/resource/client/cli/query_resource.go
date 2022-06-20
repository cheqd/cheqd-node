package cli

import (
	"context"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetResource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resource [collectionId] [id]",
		Short: "Query a resource",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			collectionId := args[0]
			id := args[1]

			params := &types.QueryGetResourceRequest{
				CollectionId: collectionId,
				Id:           id,
			}

			resp, err := queryClient.Resource(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
