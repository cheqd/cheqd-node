package cli

import (
	"context"

	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetAllResourceVersions() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-resource-versions [collectionId] [name] [resourceType] [mediaType]",
		Short: "Query all resource versions",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			collectionId := args[0]
			name := args[1]

			params := &types.QueryGetAllResourceVersionsRequest{
				CollectionId: collectionId,
				Name:         name,
			}

			resp, err := queryClient.AllResourceVersions(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
