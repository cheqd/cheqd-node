package cli

import (
	"os"

	sdkmath "cosmossdk.io/math"
	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func CmdCreateResource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create [payload-file] [resource-data-file]",
		Short: "Create a new Resource.",
		Long: `Create a new Resource within a DID Resource Collection. 
[payload-file] is JSON encoded MsgCreateResourcePayload alongside with sign inputs. 
[resource-data-file] is a path to the Resource data file.

NOTES:
1. Fee used for the transaction will ALWAYS take the fixed fee for Resource creation, REGARDLESS of what value is passed in '--fees' flag.
2. Fixed fees for Resource creation is defined based on the IANA media type of the Resource data file. These parameters can be updated using governance proposals. Currently, there are three categories of media types with different fees: 'image', 'json', and 'default' (for all other media types).
2. Payload file should contain the properties given in example below.
3. Private key provided in sign inputs is ONLY used locally to generate signature(s) and not sent to the ledger.

Example payload file:
{
    "payload": {
        "collectionId": "<did-unique-identifier>",
        "id": "<uuid>",
        "name": "<human-readable resource name>",
        "version": "<human-readable version number>",
        "resourceType": "<resource-type>",
        "alsoKnownAs": [
            {
                "uri": "did:cheqd:<namespace>:<unique-identifier>/resource/<uuid>",
                "description": "did-url"
            },
            {
                "uri": "https://example.com/alternative-uri",
                "description": "http-url"
            }
        ],
		"previousVersionResourceId": "<uuid id of the previous version>"
    },
    "signInputs": [
        {
            "verificationMethodId": "did:cheqd:<namespace>:<unique-identifier>#<key-id>",
            "privKey": "<private-key-bytes-encoded-to-base64>"
        }
    ]
}
`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// Read payload file arg
			payloadFile := args[0]

			// Read data file arg
			dataFile := args[1]

			payloadJSON, signInputs, err := didcli.ReadPayloadWithSignInputsFromFile(payloadFile)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var payload types.MsgCreateResourcePayload
			err = clientCtx.Codec.UnmarshalJSON(payloadJSON, &payload)
			if err != nil {
				return err
			}

			// Read data file
			data, err := os.ReadFile(dataFile)
			if err != nil {
				return err
			}

			// Prepare payload
			payload.Data = data

			// Populate resource id if not set
			if payload.Id == "" {
				payload.Id = uuid.NewString()
			}

			// Build identity message
			signBytes := payload.GetSignBytes()
			identitySignatures := didcli.SignWithSignInputs(signBytes, signInputs)

			msg := types.MsgCreateResource{
				Payload:    &payload,
				Signatures: identitySignatures,
			}

			// Set fee-payer if not set
			err = didcli.SetFeePayerFromSigner(&clientCtx)
			if err != nil {
				return err
			}

			return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), &msg)
		},
	}

	// add standard tx flags
	AddTxFlagsToCmd(cmd)

	// add custom / override flags
	cmd.Flags().String(flags.FlagFees, sdk.NewCoin(types.BaseMinimalDenom, sdkmath.NewInt(types.DefaultCreateResourceImageFee)).String(), "Fixed fee for Resource creation, e.g., 10000000000ncheq. Please check what the current fees by running 'cheqd-noded query params subspace resource feeparams'")

	_ = cmd.MarkFlagRequired(flags.FlagFees)
	_ = cmd.MarkFlagRequired(flags.FlagGas)
	_ = cmd.MarkFlagRequired(flags.FlagGasAdjustment)

	return cmd
}
