package cli

import (
	sdkmath "cosmossdk.io/math"
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func CmdDeactivateDidDoc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deactivate-did [payload-file] --version-id [version-id]",
		Short: "Deactivate a DID.",
		Long: `Deactivates a DID and its associated DID Document. 
[payload-file] is JSON encoded MsgDeactivateDidDocPayload alongside with sign inputs. 

NOTES:
1. Fee used for the transaction will ALWAYS take the fixed fee for DID Document deactivation, REGARDLESS of what value is passed in '--fees' flag.
2. A new DID Document version is created when deactivating a DID Document so that the operation timestamp can be recorded. Version ID is optional and is determined by the '--version-id' flag. If not provided, a random UUID will be used as version-id.
3. Payload file should be a JSON file containing the properties given in example below.
4. Private key provided in sign inputs is ONLY used locally to generate signature(s) and not sent to the ledger.

Example payload file:
{
    "payload": {
        "id": "did:cheqd:<namespace>:<unique-identifier>"
    },
    "signInputs": [
        {
            "verificationMethodId": "did:cheqd:<namespace>:<unique-identifier>#<key-id>",
            "privKey": "<private-key-bytes-encoded-to-base64>"
        }
    ]
}
`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Read payload file arg
			payloadFile := args[0]

			// Read version-id flag
			versionID, err := cmd.Flags().GetString(FlagVersionID)
			if err != nil {
				return err
			}

			if versionID != "" {
				err = utils.ValidateUUID(versionID)
				if err != nil {
					return err
				}
			} else {
				versionID = uuid.NewString()
			}

			payloadJSON, signInputs, err := ReadPayloadWithSignInputsFromFile(payloadFile)
			if err != nil {
				return err
			}

			// Build payload
			payload := types.MsgDeactivateDidDocPayload{}
			err = clientCtx.Codec.UnmarshalJSON([]byte(payloadJSON), &payload)
			if err != nil {
				return err
			}

			// Set version id from flag or random
			payload.VersionId = versionID

			// Build identity message
			signBytes := payload.GetSignBytes()
			identitySignatures := SignWithSignInputs(signBytes, signInputs)

			msg := types.MsgDeactivateDidDoc{
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

	// add standard tx flags
	AddTxFlagsToCmd(cmd)

	// add custom / override flags
	cmd.Flags().String(FlagVersionID, "", "Version ID of the DID Document")
	cmd.Flags().String(flags.FlagFees, sdk.NewCoin(types.BaseMinimalDenom, sdkmath.NewInt(types.DefaultDeactivateDidTxFee)).String(), "Fixed fee for DID deactivation, e.g., 10000000000ncheq. Please check what the current fees by running 'cheqd-noded query params subspace cheqd feeparams'")

	_ = cmd.MarkFlagRequired(flags.FlagFees)
	_ = cmd.MarkFlagRequired(flags.FlagGas)
	_ = cmd.MarkFlagRequired(flags.FlagGasAdjustment)

	return cmd
}
