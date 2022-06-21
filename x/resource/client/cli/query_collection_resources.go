package cli

import (
	"context"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetCollectionResources() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "collection-resources [collectionId]",
		Short: "Query all resource of a specific collection",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			collectionId := args[0]

			params := &types.QueryGetCollectionResourcesRequest{
				CollectionId: collectionId,
			}

			resp, err := queryClient.CollectionResources(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
