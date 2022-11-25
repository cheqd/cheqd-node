package cli

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetAllDidDocVersionsMetadata() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "all-diddoc-versions-metadata [id]",
		Short: "Query diddoc version metadata by diddoc id",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			params := &types.QueryGetAllDidDocVersionsMetadataRequest{
				Id: did,
			}

			resp, err := queryClient.AllDidDocVersionsMetadata(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
