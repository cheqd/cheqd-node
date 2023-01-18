package cli

import (
	"encoding/json"
	"fmt"

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
		Use:   "update-did [payload-file]",
		Short: "Updates a DID and its associated DID Document.",
		Long: `Update DID Document associated with a given DID. 
[payload-file] is JSON encoded DID Document alongside with sign inputs. 
Version ID is optional and is determined by the '--version-id' flag. 
If not provided, a random UUID will be used as version-id.`,
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
				Id:                   specPayload.ID,
				Controller:           specPayload.Controller,
				VerificationMethod:   verificationMethod,
				Authentication:       specPayload.Authentication,
				AssertionMethod:      specPayload.AssertionMethod,
				CapabilityInvocation: specPayload.CapabilityInvocation,
				CapabilityDelegation: specPayload.CapabilityDelegation,
				KeyAgreement:         specPayload.KeyAgreement,
				Service:              service,
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
	cmd.Flags().String(flags.FlagFees, sdk.NewCoin(types.BaseMinimalDenom, sdk.NewInt(types.DefaultUpdateDidTxFee)).String(), "Fees to pay along with transaction; eg: 25000000000ncheq")
	cmd.Flags().String(flags.FlagGas, flags.GasFlagAuto, fmt.Sprintf("gas limit to set per-transaction; set to %q to calculate sufficient gas automatically. Note: %q option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of %q. (default %d)",
		flags.GasFlagAuto, flags.GasFlagAuto, flags.FlagFees, flags.DefaultGasLimit))

	_ = cmd.MarkFlagRequired(flags.FlagFees)
	_ = cmd.MarkFlagRequired(flags.FlagGas)

	return cmd
}
