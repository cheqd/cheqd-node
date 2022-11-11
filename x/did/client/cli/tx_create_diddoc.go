package cli

import (
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/spf13/cobra"
)

func CmdCreateDidDoc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-did [payload-json] [ver-method-id-1] [priv-key-1] [ver-method-id-N] [priv-key-N] ...",
		Short: "Creates a new DID and associated DIDDoc.",
		Long: "Creates a new DID and associated DIDDoc. " +
			"[payload-json] is JSON encoded MsgCreateDidPayload. " +
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

			payloadJson, signInputs, err := GetPayloadAndSignInputs(clientCtx, args)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var payload types.MsgCreateDidDocPayload
			err = clientCtx.Codec.UnmarshalJSON([]byte(payloadJson), &payload)
			if err != nil {
				return err
			}

			// Build identity message
			signBytes := payload.GetSignBytes()
			identitySignatures := SignWithSignInputs(signBytes, signInputs)

			msg := types.MsgCreateDidDoc{
				Payload:    &payload,
				Signatures: identitySignatures,
			}

			// Set fee-payer if not set
			err = SetFeePayerFromSigner(&clientCtx)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
