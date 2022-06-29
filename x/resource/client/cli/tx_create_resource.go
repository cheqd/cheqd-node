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

const (
	FlagCollectionId = "collection-id"
	FlagResourceId   = "resource-id"
	FlagResourceName = "resource-name"
	FlagResourceType = "resource-type"
	FlagResourceFile = "resource-file"
)

func CmdCreateResource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-resource [ver-method-id-1] [priv-key-1] [ver-method-id-N] [priv-key-N] ...",
		Short: "Creates a new Resource.",
		Long: "Creates a new Resource. " +
			"[ver-method-id-N] is the DID fragment that points to the public part of the key in the ledger for the signature N." +
			"[priv-key-N] is base base64 encoded ed25519 private key for signature N." +
			"If 'interactive' value is used for a key, the key will be read interactively. " +
			"Prefer interactive mode, use inline mode only for tests.",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			collectionId, err := cmd.Flags().GetString(FlagCollectionId)
			if err != nil {
				return err
			}

			resourceId, err := cmd.Flags().GetString(FlagResourceId)
			if err != nil {
				return err
			}

			resourceName, err := cmd.Flags().GetString(FlagResourceName)
			if err != nil {
				return err
			}

			resourceType, err := cmd.Flags().GetString(FlagResourceType)
			if err != nil {
				return err
			}

			resourceFile, err := cmd.Flags().GetString(FlagResourceFile)
			if err != nil {
				return err
			}

			data, err := ioutil.ReadFile(resourceFile)
			if err != nil {
				return err
			}

			// Prepare payload
			payload := types.MsgCreateResourcePayload{
				CollectionId: collectionId,
				Id:           resourceId,
				Name:         resourceName,
				ResourceType: resourceType,
				Data:         data,
			}

			// Read signatures
			signInputs, err := cheqdcli.GetSignInputs(clientCtx, args)
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

	cmd.Flags().String(FlagCollectionId, "", "Collection ID")
	err := cobra.MarkFlagRequired(cmd.Flags(), FlagCollectionId)
	panicIfErr(err)

	cmd.Flags().String(FlagResourceId, "", "Resource ID")
	err = cobra.MarkFlagRequired(cmd.Flags(), FlagResourceId)
	panicIfErr(err)

	cmd.Flags().String(FlagResourceName, "", "Resource name")
	err = cobra.MarkFlagRequired(cmd.Flags(), FlagResourceName)
	panicIfErr(err)

	cmd.Flags().String(FlagResourceType, "", "Resource type")
	err = cobra.MarkFlagRequired(cmd.Flags(), FlagResourceType)
	panicIfErr(err)

	cmd.Flags().String(FlagResourceFile, "", "Resource file")
	err = cobra.MarkFlagRequired(cmd.Flags(), FlagResourceFile)
	panicIfErr(err)

	return cmd
}
