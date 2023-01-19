package cli

import (
	"fmt"
	"os"

	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	didutils "github.com/cheqd/cheqd-node/x/did/utils"
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
		Use:   "create [resource-payload-file] [resource-data-file]",
		Short: "Create a new Resource.",
		Long: `Create a new Resource within a DID Resource Collection. 
[resource-payload-file] is JSON encoded MsgCreateResourcePayload alongside with sign inputs. 
[resource-data-file] is a path to the Resource data file.`,
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

			// Read resource-id flag
			resourceID, err := cmd.Flags().GetString(FlagResourceID)
			if err != nil {
				return err
			}

			if resourceID != "" {
				err = didutils.ValidateUUID(resourceID)
				if err != nil {
					return err
				}
			} else {
				resourceID = uuid.NewString()
			}

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
			payload = types.MsgCreateResourcePayload{
				CollectionId: payload.CollectionId,
				Id:           resourceID,
				Name:         payload.Name,
				Version:      payload.Version,
				ResourceType: payload.ResourceType,
				AlsoKnownAs:  payload.AlsoKnownAs,
				Data:         data,
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
	didcli.AddTxFlagsToCmd(cmd)

	// add custom / override flags
	cmd.Flags().String(FlagResourceID, "", "The Resource ID. If not set, a random UUID will be generated.")
	cmd.Flags().String(flags.FlagFees, sdk.NewCoin(types.BaseMinimalDenom, sdk.NewInt(types.DefaultCreateResourceImageFee)).String(), "Fees to pay along with transaction; eg: 10000000000ncheq")
	cmd.Flags().String(flags.FlagGas, flags.GasFlagAuto, fmt.Sprintf("gas limit to set per-transaction; set to %q to calculate sufficient gas automatically. Note: %q option doesn't always report accurate results. Set a valid coin value to adjust the result. Can be used instead of %q. (default %d)",
		flags.GasFlagAuto, flags.GasFlagAuto, flags.FlagFees, flags.DefaultGasLimit))

	_ = cmd.MarkFlagRequired(flags.FlagFees)
	_ = cmd.MarkFlagRequired(flags.FlagGas)

	return cmd
}
