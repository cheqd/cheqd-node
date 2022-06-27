package cli

import (
	"io/ioutil"

	cheqdcli "github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateResource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-resource [collection-id] [id] [name] [resource-type] [mime-type] [file-path] [ver-method-id-1] [priv-key-1] [ver-method-id-N] [priv-key-N] ...",
		Short: "Creates a new Resource.",
		Long: "Creates a new Resource. " +
			"[ver-method-id-N] is the DID fragment that points to the public part of the key in the ledger for the signature N." +
			"[priv-key-1] is base base64 encoded ed25519 private key for signature N." +
			"If 'interactive' value is used for a key, the key will be read interactively. " +
			"Prefer interactive mode, use inline mode only for tests.",
		Args: cobra.MinimumNArgs(6),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionId := args[0]

			id := args[1]
			name := args[2]
			resourceType := args[3]
			mediaType := args[4]

			data, err := ioutil.ReadFile(args[5])
			if err != nil {
				return err
			}

			// Prepare payload
			payload := types.MsgCreateResourcePayload{
				CollectionId: collectionId,
				Id:           id,
				Name:         name,
				ResourceType: resourceType,
				MediaType:    mediaType,
				Data:         data,
			}

			// Read signatures
			signInputs, err := cheqdcli.GetSignInputs(clientCtx, args[6:])
			if err != nil {
				return err
			}

			// Build identity message
			signBytes := payload.GetSignBytes()
			identitySignatures := cheqdcli.SignWithSignInputs(signBytes, signInputs)

			msg := types.MsgCreateResource{
				Payload:    &payload,
				Signatures: identitySignatures,
			}

			// Set fee-payer if not set
			err = cheqdcli.SetFeePayerFromSigner(&clientCtx)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
