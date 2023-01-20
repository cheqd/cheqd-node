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
		Use:   "specific-resource [collection-id] [resource-id]",
		Short: "Query a specific resource",
		Long: `Query a specific resource by Collection ID and Resource ID.
		
		Collection ID is the UNIQUE IDENTIFIER part of the DID the resource is linked to.
		Example: c82f2b02-bdab-4dd7-b833-3e143745d612, wGHEXrZvJxR8vw5P3UWH1j, etc.

		Resource ID is the UUID of the specific resource.
		Example: 6e8bc430-9c3a-11d9-9669-0800200c9a66, 6e8bc430-9c3a-11d9-9669-0800200c9a67, etc.`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			collectionID := args[0]
			id := args[1]

			params := &types.QueryResourceRequest{
				CollectionId: collectionID,
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
