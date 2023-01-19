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

func CmdCreateDidDoc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-did [payload-file]",
		Short: "Create a new DID and associated DID Document.",
		Long: `Creates a new DID and associated DID Document. 
[payload-file] is JSON encoded DID Document alongside with sign inputs.
Version ID is optional and is determined by the '--version-id' flag. 
If not provided, a random UUID will be used as version-id.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			payloadFile := args[0]
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

			// Construct MsgCreateDidDocPayload
			payload := types.MsgCreateDidDocPayload{
				Id:                   specPayload.ID,
				Controller:           specPayload.Controller,
				VerificationMethod:   verificationMethod,
				Authentication:       specPayload.Authentication,
				AssertionMethod:      specPayload.AssertionMethod,
				CapabilityInvocation: specPayload.CapabilityInvocation,
				CapabilityDelegation: specPayload.CapabilityDelegation,
				KeyAgreement:         specPayload.KeyAgreement,
				Service:              service,
				VersionId:            versionID,
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

	// add standard tx flags
	AddTxFlagsToCmd(cmd)

	// add custom / override flags
	cmd.Flags().String(FlagVersionID, "", "Version ID of the DID Document")
	cmd.Flags().String(flags.FlagFees, sdk.NewCoin(types.BaseMinimalDenom, sdk.NewInt(types.DefaultCreateDidTxFee)).String(), "Fees to pay along with transaction; eg: 50000000000ncheq")

	_ = cmd.MarkFlagRequired(flags.FlagFees)
	_ = cmd.MarkFlagRequired(flags.FlagGas)
	_ = cmd.MarkFlagRequired(flags.FlagGasAdjustment)

	return cmd
}
