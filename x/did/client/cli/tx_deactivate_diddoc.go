package cli

import (
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func CmdDeactivateDidDoc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivate-diddoc [id] [previous-version-id] [ver-method-id-1] [priv-key-1] [ver-method-id-N] [priv-key-N] ...",
		Short: "Deactivates a DIDDoc.",
		Long: "Deactivates a DIDDoc. " +
			"[id] is DID Document id. " +
			"[ver-method-id-N] is the DID fragment that points to the public part of the key in the ledger for the signature N." +
			"[priv-key-1] is base base64 encoded ed25519 private key for signature N." +
			"If 'interactive' value is used for a key, the key will be read interactively. " +
			"Prefer interactive mode, use inline mode only for tests.",
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			did := args[0]
			previousVersionId := args[1]

			signInputs, err := GetSignInputs(clientCtx, args[2:])
			if err != nil {
				return err
			}

			// Build payload
			payload := &types.MsgDeactivateDidDocPayload{
				Id:                did,
				PreviousVersionId: previousVersionId,
				VersionId:         uuid.NewString(),
			}

			// Build identity message
			signBytes := payload.GetSignBytes()
			identitySignatures := SignWithSignInputs(signBytes, signInputs)

			msg := types.MsgDeactivateDidDoc{
				Payload:    payload,
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
