package cli

import (
	"encoding/json"

	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cheqd/cheqd-node/x/did/utils"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func CmdUpdateDidDoc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update-did [payload-file] --version-id [version-id]",
		Short: "Updates a DID and its associated DID Document.",
		Long: `Update DID Document associated with a given DID. 
[payload-file] is JSON encoded DID Document alongside with sign inputs. 
Version ID is optional and is determined by the '--version-id' flag. 
If not provided, a random UUID will be used as version-id.

NOTES:
1. Fee used for the transaction will ALWAYS take the fixed fee for DID Document update, REGARDLESS of what value is passed in '--fees' flag.
2. DID update operations require the FULL new DID Document to be provided. Specifying just the changes/diff is not supported.
3. Payload file should be a JSON file containing properties specified in the DID Core Specification. Rules from DID Core spec are followed on which properties are mandatory and which ones are optional.
4. Private key provided in sign inputs is ONLY used locally to generate signature(s) and not sent to the ledger.

Example payload file:
{
    "payload": {
        "context": [ "https://www.w3.org/ns/did/v1" ],
        "id": "did:cheqd:<namespace>:<unique-identifier>",
        "controller": [
            "did:cheqd:<namespace>:<unique-identifier>"
        ],
        "authentication": [
            "did:cheqd:<namespace>:<unique-identifier>#<key-id>"
        ],
        "assertionMethod": [],
        "capabilityInvocation": [],
        "capabilityDelegation": [],
        "keyAgreement": [],
        "alsoKnownAs": [],
        "verificationMethod": [
            {
                "id": "did:cheqd:<namespace>:<unique-identifier>#<key-id>",
                "type": "<verification-method-type>",
                "controller": "did:cheqd:<namespace>:<unique-identifier>",
                "publicKeyMultibase": "<public-key>"
            }
        ],
        "service": [
			{
                "id": "did:cheqd:<namespace>:<unique-identifier>#<service-id>",
                "type": "<service-type>",
                "serviceEndpoint": [
                    "<service-endpoint>"
                ]
            }
		]
    },
	"signInputs": [
        {
            "verificationMethodId": "did:cheqd:<namespace>:<unique-identifier>#<key-id>",
            "privKey": "<private-key-bytes-encoded-to-base64>"
        }
    ]
}`,
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

			// Unmarshal spec-compliant payload
			var specPayload DIDDocument
			err = json.Unmarshal([]byte(payloadJSON), &specPayload)
			if err != nil {
				return err
			}

			// Validate spec-compliant payload & get verification methods
			verificationMethod, service, err := GetFromSpecCompliantPayload(specPayload)
			if err != nil {
				return err
			}

			// Construct MsgUpdateDidDocPayload
			payload := types.MsgUpdateDidDocPayload{
				Context:              specPayload.Context,
				Id:                   specPayload.ID,
				Controller:           specPayload.Controller,
				VerificationMethod:   verificationMethod,
				Authentication:       specPayload.Authentication,
				AssertionMethod:      specPayload.AssertionMethod,
				CapabilityInvocation: specPayload.CapabilityInvocation,
				CapabilityDelegation: specPayload.CapabilityDelegation,
				KeyAgreement:         specPayload.KeyAgreement,
				Service:              service,
				AlsoKnownAs:          specPayload.AlsoKnownAs,
				VersionId:            versionID, // Set version id, from flag or random
			}

			// Build identity message
			signBytes := payload.GetSignBytes()
			identitySignatures := SignWithSignInputs(signBytes, signInputs)

			msg := types.MsgUpdateDidDoc{
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
	cmd.Flags().String(flags.FlagFees, sdk.NewCoin(types.BaseMinimalDenom, sdk.NewInt(types.DefaultUpdateDidTxFee)).String(), "Fixed fee for DID update, e.g., 25000000000ncheq. Please check what the current fees by running 'cheqd-noded query params subspace cheqd feeparams'")

	_ = cmd.MarkFlagRequired(flags.FlagFees)
	_ = cmd.MarkFlagRequired(flags.FlagGas)
	_ = cmd.MarkFlagRequired(flags.FlagGasAdjustment)

	return cmd
}
