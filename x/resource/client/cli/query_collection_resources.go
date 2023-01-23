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
		Use:   "collection-metadata [collection-id] [resource-id]",
		Short: "Query metadata for an entire Collection",
		Long: `Query metadata for an entire Collection by Collection ID. This will return the metadata for all Resources in the Collection.
		
		Collection ID is the UNIQUE IDENTIFIER part of the DID the resource is linked to.
		Example: c82f2b02-bdab-4dd7-b833-3e143745d612, wGHEXrZvJxR8vw5P3UWH1j, etc.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			collectionID := args[0]

			params := &types.QueryCollectionResourcesRequest{
				CollectionId: collectionID,
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
