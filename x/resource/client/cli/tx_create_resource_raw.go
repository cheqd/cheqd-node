package cli

import (
	cheqdcli "github.com/cheqd/cheqd-node/x/cheqd/client/cli"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateResourceRaw() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-resource-raw [payload-json] [ver-method-id-1] [priv-key-1] [ver-method-id-N] [priv-key-N] ...",
		Short: "Creates a new Resource using raw payload. For testing purposes.",
		Long: "Creates a new Resource. " +
			"[payload-json] is JSON encoded MsgCreateResourcePayload. " +
			"[ver-method-id-N] is the DID fragment that points to the public part of the key in the ledger for the signature N." +
			"[priv-key-1] is base base64 encoded ed25519 private key for signature N." +
			"If 'interactive' value is used for a key, the key will be read interactively. " +
			"Prefer interactive mode, use inline mode only for tests.",
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			payloadJson, signInputs, err := cheqdcli.GetPayloadAndSignInputs(clientCtx, args)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var payload types.MsgCreateResourcePayload
			err = clientCtx.Codec.UnmarshalJSON([]byte(payloadJson), &payload)
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
