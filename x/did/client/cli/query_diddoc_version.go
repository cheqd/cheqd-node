package cli

import (
	"context"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/spf13/cobra"
)

func CmdGetDidDocVersion() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "diddoc-version [id] [version-id]",
		Short: "Query a specific diddoc version",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx := client.GetClientContextFromCmd(cmd)

			queryClient := types.NewQueryClient(clientCtx)

			did := args[0]
			versionId := args[1]
			params := &types.QueryGetDidDocVersionRequest{
				Id:      did,
				Version: versionId,
			}

			resp, err := queryClient.DidDocVersion(context.Background(), params)
			if err != nil {
				return err
			}

			return clientCtx.PrintProto(resp)
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}
