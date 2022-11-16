package cli

import (
	"encoding/json"
	"os"

	didcli "github.com/cheqd/cheqd-node/x/did/client/cli"
	"github.com/cheqd/cheqd-node/x/resource/types"
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/google/uuid"
	"github.com/spf13/cobra"
)

type CreateResourceOptions struct {
	CollectionId    string                  `json:"collection_id"`
	ResourceId      string                  `json:"resource_id"`
	ResourceName    string                  `json:"resource_name"`
	ResourceVersion string                  `json:"resource_version"`
	ResourceType    string                  `json:"resource_type"`
	ResourceFile    string                  `json:"resource_file"`
	AlsoKnownAs     []*types.AlternativeUri `json:"also_known_as"`
}

func CmdCreateResource() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create-resource [payload-file]",
		Short: "Creates a new Resource.",
		Long: "Creates a new Resource. " +
			"[payload-file] is JSON encoded MsgCreateDidDocPayload alongside with sign inputs.",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			clientCtx, err := client.GetClientTxContext(cmd)
			if err != nil {
				return err
			}

			payloadFile := args[0]

			payloadJson, signInputs, err := didcli.ReadPayloadWithSignInputsFromFile(payloadFile)
			if err != nil {
				return err
			}

			var options CreateResourceOptions
			err = json.Unmarshal(payloadJson, &options)
			if err != nil {
				return err
			}

			data, err := os.ReadFile(options.ResourceFile)
			if err != nil {
				return err
			}

			// Prepare payload
			payload := types.MsgCreateResourcePayload{
				CollectionId: options.CollectionId,
				Id:           options.ResourceId,
				Name:         options.ResourceName,
				Version:      options.ResourceVersion,
				ResourceType: options.ResourceType,
				AlsoKnownAs:  options.AlsoKnownAs,
				Data:         data,
			}

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

	flags.AddTxFlagsToCmd(cmd)

	return cmd
}
