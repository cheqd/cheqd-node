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
		Use:   "did-metadata [id]",
		Short: "Query all versions metadata for a DID",
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
