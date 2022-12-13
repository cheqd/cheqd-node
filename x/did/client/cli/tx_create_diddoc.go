package cli

import (
	"github.com/cheqd/cheqd-node/x/did/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

func CmdCreateDidDoc() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-did [payload-file]",
		Short: "Create a new DID and associated DID Document.",
		Long: "Creates a new DID and associated DID Document. " +
			"[payload-file] is JSON encoded MsgCreateDidDocPayload alongside with sign inputs.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			// 1. Add utility-fee as 2nd arg
			// 2. Short
			// 3. Long
			// 4. Show example to query the actual params `cheqd-noded query params subspace cheqd`

			payloadFile := args[0]

			payloadJson, signInputs, err := ReadPayloadWithSignInputsFromFile(payloadFile)
			if err != nil {
				return err
			}

			// Unmarshal payload
			var payload types.MsgCreateDidDocPayload
			err = clientCtx.Codec.UnmarshalJSON([]byte(payloadJson), &payload)
			if err != nil {
				return err
			}

			// Check for versionId
			if payload.VersionId == "" {
				payload.VersionId = uuid.NewString()
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
